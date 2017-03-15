package bc

import "errors"

var (
	errPosition              = errors.New("invalid source or destination position")
	errEntryType             = errors.New("invalid entry type")
	errBadTimeRange          = errors.New("bad time range")
	errEmptyResults          = errors.New("transaction has no results")
	errMismatchedAssetID     = errors.New("mismatched asset id")
	errMismatchedBlock       = errors.New("mismatched block")
	errMismatchedMerkleRoot  = errors.New("mismatched merkle root")
	errMismatchedReference   = errors.New("mismatched reference")
	errMismatchedValue       = errors.New("mismatched value")
	errMisorderedBlockHeight = errors.New("misordered block height")
	errMisorderedBlockTime   = errors.New("misordered block time")
	errNoPrevBlock           = errors.New("no previous block")
	errNoSource              = errors.New("no source for value")
	errNonemptyExtHash       = errors.New("non-empty extension hash")
	errOverflow              = errors.New("arithmetic overflow/underflow")
	errTxVersion             = errors.New("invalid transaction version")
	errUnbalanced            = errors.New("unbalanced")
	errUntimelyTransaction   = errors.New("block timestamp outside transaction time range")
	errVersionRegression     = errors.New("version regression")
	errWrongBlockchain       = errors.New("wrong blockchain")
	errZeroTime              = errors.New("timerange has one or two bounds set to zero")
)

type validationState struct {
	blockVersion   uint64
	initialBlockID Hash

	currentTx *TxEntries

	// Set this to the ID of an entry whenever validating an entry
	currentEntryID Hash

	// Must be defined when validating a valueSource
	sourcePosition uint64

	// Must be defined when validating a valueDestination
	destPosition uint64

	// The block timestamp
	timestampMS       uint64
	prevBlockHeader   *BlockHeaderEntry
	prevBlockHeaderID Hash
	blockTxs          []*TxEntries
	// xxx reachable entries?
}

type blockVMContext struct {
	prog  Program
	args  [][]byte
	block *BlockEntries
}

func (b *blockVMContext) VMVersion() uint64   { return b.prog.VMVersion }
func (b *blockVMContext) Code() []byte        { return b.prog.Code }
func (b *blockVMContext) Arguments() [][]byte { return b.args }

func (b *blockVMContext) BlockHash() ([]byte, error)   { return b.block.ID[:], nil }
func (b *blockVMContext) BlockTimeMS() (uint64, error) { return b.block.body.TimestampMS }

func (b *blockVMContext) NextConsensusProgram() ([]byte, error) {
	return b.block.body.NextConsensusProgram, nil
}

func (b *blockVMContext) TxVersion() (uint64, bool)      { return 0, false }
func (b *blockVMContext) TxSigHash() ([]byte, error)     { return nil, ErrContext }
func (b *blockVMContext) NumResults() (uint64, error)    { return 0, ErrContext }
func (b *blockVMContext) AssetID() ([]byte, error)       { return nil, ErrContext }
func (b *blockVMContext) Amount() (uint64, error)        { return 0, ErrContext }
func (b *blockVMContext) MinTimeMS() (uint64, error)     { return 0, ErrContext }
func (b *blockVMContext) MaxTimeMS() (uint64, error)     { return 0, ErrContext }
func (b *blockVMContext) EntryData() ([]byte, error)     { return nil, ErrContext } // xxx ?
func (b *blockVMContext) TxData() ([]byte, error)        { return nil, ErrContext }
func (b *blockVMContext) DestPos() (uint64, error)       { return 0, ErrContext }
func (b *blockVMContext) AnchorID() ([]byte, error)      { return nil, ErrContext }
func (b *blockVMContext) SpentOutputID() ([]byte, error) { return nil, ErrContext }

func (b *blockVMContext) CheckOutput(uint64, []byte, uint64, []byte, uint64, []byte) (bool, error) {
	return false, ErrContext
}

func newBlockVMContext(blockEntries *BlockEntries, prog []byte, args [][]byte) *blockVMContext {
	return &blockVMContext{}
}

type txVMContext struct {
}

func newTxVMContext(txEntries *TxEntries, entry Entry, prog Program, args [][]byte) *txVMContext {
	return &txVMContext{}
}
