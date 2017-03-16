package bc

import (
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

func (bh *BlockHeaderEntry) CheckValid(state *validationState) error {
	if state.prevBlockHeader == nil {
		if bh.body.Height != 1 {
			return errors.WithDetailf(errNoPrevBlock, "height %d", bh.body.Height)
		}
	} else {
		if bh.body.Version < state.prevBlockHeader.body.Version {
			return errors.WithDetailf(errVersionRegression, "previous block verson %d, current block version %d", state.prevBlockHeader.body.Version, bh.body.Version)
		}

		if bh.body.Height != state.prevBlockHeader.body.Height+1 {
			return errors.WithDetailf(errMisorderedBlockHeight, "previous block height %d, current block height %d", state.prevBlockHeader.body.Height, bh.body.Height)
		}

		if state.prevBlockHeaderID != bh.body.PreviousBlockID {
			return errors.WithDetailf(errMismatchedBlock, "previous block ID %x, current block wants %x", state.prevBlockHeaderID[:], bh.body.PreviousBlockID[:])
		}

		if bh.body.TimestampMS <= state.prevBlockHeader.body.TimestampMS {
			return errors.WithDetailf(errMisorderedBlockTime, "previous block time %d, current block time %d", state.prevBlockHeader.body.TimestampMS, bh.body.TimestampMS)
		}

		if state.blockVMContext != nil {
			err := vm.Verify(state.blockVMContext)
			if err != nil {
				return errors.Wrap(err, "evaluating previous block's next consensus program")
			}
		}
	}

	for i, tx := range state.blockTxs {
		txState := *state
		txState.currentEntryID = tx.ID
		txState.currentTx = tx
		err := tx.CheckValid(&txState)
		if err != nil {
			return errors.Wrapf(err, "checking validity of transaction %d of %d", i, len(state.blockTxs))
		}
	}

	txRoot, err := CalcMerkleRoot(state.blockTxs)
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
