package validation

import (
	"bytes"
	"context"
	"encoding/hex"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"

	"chain/errors"
	"chain/protocol"
	"chain/protocol/bc"
	"chain/protocol/vm"
	"chain/protocol/vmutil"
)

// Errors returned by the block validation functions.
var (
	ErrBadPrevHash  = errors.New("invalid previous block hash")
	ErrBadHeight    = errors.New("invalid block height")
	ErrBadTimestamp = errors.New("invalid block timestamp")
	ErrBadScript    = errors.New("unspendable block script")
	ErrBadSig       = errors.New("invalid signature script")
	ErrBadTxRoot    = errors.New("invalid transaction merkle root")
	ErrBadStateRoot = errors.New("invalid state merkle root")
)

// ValidateBlock performs the "validate block" procedure from the spec,
// yielding a new state (recorded in the 'snapshot' argument).
// See $CHAIN/protocol/doc/spec/validation.md#validate-block.
// Note that it does not execute prevBlock's consensus program.
// (See ValidateBlockForAccept for that.)
func ValidateBlock(ctx context.Context, snapshot *protocol.Snapshot, initialBlockHash bc.Hash, prevBlock, block *bc.BlockEntries) error {

	var g errgroup.Group
	// Do all of the unparallelizable work, plus validating the block
	// header in one goroutine.
	g.Go(func() error {
		var prev *bc.BlockHeaderEntry
		if prevBlock != nil {
			prev = prevBlock.BlockHeaderEntry
		}
		err := validateBlockHeader(prev, block)
		if err != nil {
			return err
		}
		snapshot.PruneNonces(block.TimestampMS())

		// TODO: Check that other block headers are valid.
		// TODO(erykwalder): consider writing to a copy of the state tree
		// of the one provided and make the caller call ApplyBlock as well
		for _, tx := range block.Transactions {
			err = ConfirmTx(snapshot, initialBlockHash, block.Version(), block.TimestampMS(), tx)
			if err != nil {
				return err
			}
			err = ApplyTx(snapshot, tx)
			if err != nil {
				return err
			}
		}
		if block.AssetsRoot() != snapshot.Tree.RootHash() {
			return ErrBadStateRoot
		}
		return nil
	})

	// Distribute checking well-formedness of the transactions across
	// GOMAXPROCS goroutines.
	ch := make(chan *bc.TxEntries, len(block.Transactions))
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		g.Go(func() error {
			for tx := range ch {
				if err := validateTx(tx); err != nil {
					return err
				}
			}
			return nil
		})
	}
	for _, tx := range block.Transactions {
		ch <- tx
	}
	close(ch)
	return g.Wait()
}

// ApplyBlock applies the transactions in the block to the state tree.
func ApplyBlock(snapshot *protocol.Snapshot, block *bc.BlockEntries) error {
	snapshot.PruneNonces(block.TimestampMS())
	for _, tx := range block.Transactions {
		err := ApplyTx(snapshot, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateBlockHeader(prev *bc.BlockHeaderEntry, block *bc.BlockEntries) error {
	if prev == nil && block.Height() != 1 {
		return ErrBadHeight
	}
	if prev != nil {
		prevHash := bc.EntryID(prev)
		if !bytes.Equal(block.PreviousBlockID().Bytes(), prevHash[:]) {
			return ErrBadPrevHash
		}
		if block.Height() != prev.Height()+1 {
			return ErrBadHeight
		}
		if block.TimestampMS() < prev.TimestampMS() {
			return ErrBadTimestamp
		}
	}

	txMerkleRoot, err := CalcMerkleRoot(block.Transactions)
	if err != nil {
		return errors.Wrap(err, "calculating tx merkle root")
	}

	// can be modified to allow soft fork
	if block.TransactionsRoot() != txMerkleRoot {
		return ErrBadTxRoot
	}

	if vmutil.IsUnspendable(block.NextConsensusProgram()) {
		return ErrBadScript
	}

	return nil
}
