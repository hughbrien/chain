package bc

import (
	"context"

	"chain/errors"
	"chain/protocol/vm"
)

// BlockHeaderEntry contains the header information for a blockchain
// block. It satisfies the Entry interface.
type BlockHeaderEntry struct {
	body struct {
		Version              uint64
		Height               uint64
		PreviousBlockID      Hash
		TimestampMS          uint64
		TransactionsRoot     Hash
		AssetsRoot           Hash
		NextConsensusProgram []byte
		ExtHash              Hash
	}
	witness struct {
		Arguments [][]byte
	}
}

func (BlockHeaderEntry) Type() string          { return "blockheader" }
func (bh *BlockHeaderEntry) Body() interface{} { return bh.body }

func (BlockHeaderEntry) Ordinal() int { return -1 }

func (bh *BlockHeaderEntry) Version() uint64 {
	return bh.body.Version
}

func (bh *BlockHeaderEntry) Height() uint64 {
	return bh.body.Height
}

func (bh *BlockHeaderEntry) PreviousBlockID() Hash {
	return bh.body.PreviousBlockID
}

func (bh *BlockHeaderEntry) TimestampMS() uint64 {
	return bh.body.TimestampMS
}

func (bh *BlockHeaderEntry) TransactionsRoot() Hash {
	return bh.body.TransactionsRoot
}

func (bh *BlockHeaderEntry) AssetsRoot() Hash {
	return bh.body.AssetsRoot
}

func (bh *BlockHeaderEntry) NextConsensusProgram() []byte {
	return bh.body.NextConsensusProgram
}

func (bh *BlockHeaderEntry) Arguments() [][]byte {
	return bh.witness.Arguments
}

func (bh *BlockHeaderEntry) SetArguments(args [][]byte) {
	bh.witness.Arguments = args
}

// NewBlockHeaderEntry creates a new BlockHeaderEntry and populates
// its body.
func NewBlockHeaderEntry(version, height uint64, previousBlockID Hash, timestampMS uint64, transactionsRoot, assetsRoot Hash, nextConsensusProgram []byte) *BlockHeaderEntry {
	bh := new(BlockHeaderEntry)
	bh.body.Version = version
	bh.body.Height = height
	bh.body.PreviousBlockID = previousBlockID
	bh.body.TimestampMS = timestampMS
	bh.body.TransactionsRoot = transactionsRoot
	bh.body.AssetsRoot = assetsRoot
	bh.body.NextConsensusProgram = nextConsensusProgram
	return bh
}

func (bh *BlockHeaderEntry) CheckValid(ctx context.Context) error {
	bvInfo, _ := ctx.Value(vcBlockValidationInfo).(*blockValidationInfo)

	prevBlockHeader := bvInfo.prevBlockHeader

	if prevBlockHeader == nil {
		if bh.body.Height != 1 {
			return errors.WithDetailf(errNoPrevBlock, "height %d", bh.body.Height)
		}
	} else {
		if bh.body.Version < prevBlockHeader.body.Version {
			return errors.WithDetailf(errVersionRegression, "previous block verson %d, current block version %d", prevBlockHeader.body.Version, bh.body.Version)
		}

		if bh.body.Height != prevBlockHeader.body.Height+1 {
			return errors.WithDetailf(errMisorderedBlockHeight, "previous block height %d, current block height %d", prevBlockHeader.body.Height, bh.body.Height)
		}

		prevBlockHeaderID := bvInfo.prevBlockHeaderID

		if prevBlockHeaderID != bh.body.PreviousBlockID {
			return errors.WithDetailf(errMismatchedBlock, "previous block ID %x, current block wants %x", prevBlockHeaderID[:], bh.body.PreviousBlockID[:])
		}

		if bh.body.TimestampMS <= prevBlockHeader.body.TimestampMS {
			return errors.WithDetailf(errMisorderedBlockTime, "previous block time %d, current block time %d", prevBlockHeader.body.TimestampMS, bh.body.TimestampMS)
		}

		blockVMContext := bvInfo.blockVMContext

		if blockVMContext != nil {
			err := vm.Verify(blockVMContext)
			if err != nil {
				return errors.Wrap(err, "evaluating previous block's next consensus program")
			}
		}
	}

	blockTxs := bvInfo.blockTxs

	txvInfo := &txValidationInfo{
		blockVersion: bh.body.Version,
		timestampMS:  bh.body.TimestampMS,
	}
	ctx = context.WithValue(ctx, vcTxValidationInfo, txvInfo)

	for i, tx := range blockTxs {
		ctx = context.WithValue(ctx, vcCurrentEntryID, tx.ID)
		ctx = context.WithValue(ctx, vcCurrentTx, tx)
		err := tx.CheckValid(ctx)
		if err != nil {
			return errors.Wrapf(err, "checking validity of transaction %d of %d", i, len(blockTxs))
		}
	}

	txRoot, err := CalcMerkleRoot(blockTxs)
	if err != nil {
		return errors.Wrap(err, "computing transaction merkle root")
	}

	if txRoot != bh.body.TransactionsRoot {
		return errors.WithDetailf(errMismatchedMerkleRoot, "computed %x, current block wants %x", txRoot[:], bh.body.TransactionsRoot[:])
	}

	if bh.body.Version == 1 && (bh.body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
