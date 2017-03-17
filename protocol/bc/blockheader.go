package bc

import "context"

type (
	BlockHeaderEntryBody struct {
		Version              uint64
		Height               uint64
		PreviousBlockID      Hash
		TimestampMS          uint64
		TransactionsRoot     Hash
		AssetsRoot           Hash
		NextConsensusProgram []byte
		ExtHash              Hash
	}

	BlockHeaderEntryWitness struct {
		Arguments [][]byte
	}

	// BlockHeaderEntry contains the header information for a blockchain
	// block. It satisfies the Entry interface.
	BlockHeaderEntry struct {
		BlockHeaderEntryBody
		BlockHeaderEntryWitness
	}
)

func (BlockHeaderEntry) Type() string          { return "blockheader" }
func (bh *BlockHeaderEntry) Body() interface{} { return bh.BlockHeaderEntryBody }

func (BlockHeaderEntry) Ordinal() int { return -1 }

// NewBlockHeaderEntry creates a new BlockHeaderEntry and populates
// its body.
func NewBlockHeaderEntry(version, height uint64, previousBlockID Hash, timestampMS uint64, transactionsRoot, assetsRoot Hash, nextConsensusProgram []byte) *BlockHeaderEntry {
	return &BlockHeaderEntry{
		BlockHeaderEntryBody: BlockHeaderEntryBody{
			Version:              version,
			Height:               height,
			PreviousBlockID:      previousBlockID,
			TimestampMS:          timestampMS,
			TransactionsRoot:     transactionsRoot,
			AssetsRoot:           assetsRoot,
			NextConsensusProgram: nextConsensusProgram,
		},
	}
}

// CheckValid does only part of the work of validating a block. The
// rest is handled in ValidateBlock, which calls this.
func (bh *BlockHeaderEntry) CheckValid(ctx context.Context) error {
	if bh.Version == 1 && (bh.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}
	return nil
}
