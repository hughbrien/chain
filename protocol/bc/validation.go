package bc

import "errors"

var (
	errBadTimeRange          = errors.New("bad time range")
	errContext               = errors.New("wrong context")
	errEmptyResults          = errors.New("transaction has no results")
	errEntryType             = errors.New("invalid entry type")
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
	errPosition              = errors.New("invalid source or destination position")
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

	blockVMContext *blockVMContext
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
func (b *blockVMContext) BlockTimeMS() (uint64, error) { return b.block.body.TimestampMS, nil }

func (b *blockVMContext) NextConsensusProgram() ([]byte, error) {
	return b.block.body.NextConsensusProgram, nil
}

func (b *blockVMContext) TxVersion() (uint64, bool)      { return 0, false }
func (b *blockVMContext) TxSigHash() ([]byte, error)     { return nil, errContext }
func (b *blockVMContext) NumResults() (uint64, error)    { return 0, errContext }
func (b *blockVMContext) AssetID() ([]byte, error)       { return nil, errContext }
func (b *blockVMContext) Amount() (uint64, error)        { return 0, errContext }
func (b *blockVMContext) MinTimeMS() (uint64, error)     { return 0, errContext }
func (b *blockVMContext) MaxTimeMS() (uint64, error)     { return 0, errContext }
func (b *blockVMContext) EntryData() ([]byte, error)     { return nil, errContext } // xxx ?
func (b *blockVMContext) TxData() ([]byte, error)        { return nil, errContext }
func (b *blockVMContext) DestPos() (uint64, error)       { return 0, errContext }
func (b *blockVMContext) AnchorID() ([]byte, error)      { return nil, errContext }
func (b *blockVMContext) SpentOutputID() ([]byte, error) { return nil, errContext }

func (b *blockVMContext) CheckOutput(uint64, []byte, uint64, []byte, uint64, []byte) (bool, error) {
	return false, errContext
}

func newBlockVMContext(block *BlockEntries, prog []byte, args [][]byte) *blockVMContext {
	return &blockVMContext{
		prog: Program{
			VMVersion: 1,
			Code:      prog,
		},
		args:  args,
		block: block,
	}
}

type txVMContext struct {
	prog  Program
	args  [][]byte
	tx    *TxEntries
	entry Entry
}

func newTxVMContext(tx *TxEntries, entry Entry, prog Program, args [][]byte) *txVMContext {
	return &txVMContext{
		prog:  prog,
		args:  args,
		tx:    tx,
		entry: entry,
	}
}

func (t *txVMContext) VMVersion() uint64   { return t.prog.VMVersion }
func (t *txVMContext) Code() []byte        { return t.prog.Code }
func (t *txVMContext) Arguments() [][]byte { return t.args }

func (t *txVMContext) BlockHash() ([]byte, error)   { return nil, errContext }
func (t *txVMContext) BlockTimeMS() (uint64, error) { return 0, errContext }

func (t *txVMContext) NextConsensusProgram() ([]byte, error) { return nil, errContext }

func (t *txVMContext) TxVersion() (uint64, bool) { return t.tx.body.Version, true }

func (t *txVMContext) TxSigHash() ([]byte, error) {
	ord := t.entry.Ordinal()
	if ord < 0 {
		return nil, errContext
	}
	h := t.tx.SigHash(uint32(ord))
	return h[:], nil
}

func (t *txVMContext) NumResults() (uint64, error) { return uint64(len(t.tx.Results)), nil }

func (t *txVMContext) AssetID() ([]byte, error) {
	switch inp := t.entry.(type) {
	case *Nonce:
		if iss, ok := inp.Anchored.(*Issuance); ok {
			return iss.body.Value.AssetID[:], nil
		}
		return nil, errContext

	case *Issuance:
		return inp.body.Value.AssetID[:], nil

	case *Spend:
		return inp.SpentOutput.body.Source.Value.AssetID[:], nil
	}

	return nil, errContext
}

func (t *txVMContext) Amount() (uint64, error) {
	switch inp := t.entry.(type) {
	case *Nonce:
		if iss, ok := inp.Anchored.(*Issuance); ok {
			return iss.body.Value.Amount, nil
		}
		return 0, errContext

	case *Issuance:
		return inp.body.Value.Amount, nil

	case *Spend:
		return inp.SpentOutput.body.Source.Value.Amount, nil
	}

	return 0, errContext
}

func (t *txVMContext) MinTimeMS() (uint64, error) { return t.tx.body.MinTimeMS, nil }
func (t *txVMContext) MaxTimeMS() (uint64, error) { return t.tx.body.MaxTimeMS, nil }

func (t *txVMContext) EntryData() ([]byte, error) {
	switch inp := t.entry.(type) {
	case *Issuance:
		return inp.body.Data[:], nil

	case *Spend:
		return inp.body.Data[:], nil

	case *Output:
		return inp.body.Data[:], nil

	case *Retirement:
		return inp.body.Data[:], nil
	}

	return nil, errContext
}

func (t *txVMContext) TxData() ([]byte, error) { return t.tx.body.Data[:], nil }

func (t *txVMContext) DestPos() (uint64, error) {
	switch inp := t.entry.(type) {
	case *Issuance:
		return inp.witness.Destination.Position, nil

	case *Spend:
		return inp.witness.Destination.Position, nil
	}

	return 0, errContext
}

func (t *txVMContext) AnchorID() ([]byte, error) {
	if inp, ok := t.entry.(*Issuance); ok {
		return inp.body.Anchor[:], nil
	}
	return nil, errContext
}

func (t *txVMContext) SpentOutputID() ([]byte, error) {
	if inp, ok := t.entry.(*Spend); ok {
		return inp.body.SpentOutput[:], nil
	}
	return nil, errContext
}

func (t *txVMContext) CheckOutput(uint64, []byte, uint64, []byte, uint64, []byte) (bool, error) {
	return false, errContext
}
