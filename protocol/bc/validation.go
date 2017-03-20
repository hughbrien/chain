package bc

import (
	"bytes"
	"errors"

	"chain/protocol/vm"
)

var (
	errBadTimeRange          = errors.New("bad time range")
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

// keys for storing/retrieving validation values to/from context objects
const (
	vcCurrentEntryID = 1 + iota // Entry
	vcCurrentTx                 // *TxEntries
	vcSourcePos                 // uint64
	vcDestPos                   // uint64
	vcInitialBlockID            // Hash
)

type BlockVMContext struct {
	prog  Program
	args  [][]byte
	block *BlockEntries
}

func (b *BlockVMContext) VMVersion() uint64   { return b.prog.VMVersion }
func (b *BlockVMContext) Code() []byte        { return b.prog.Code }
func (b *BlockVMContext) Arguments() [][]byte { return b.args }

func (b *BlockVMContext) BlockHash() ([]byte, error)   { return b.block.ID[:], nil }
func (b *BlockVMContext) BlockTimeMS() (uint64, error) { return b.block.Body.TimestampMS, nil }

func (b *BlockVMContext) NextConsensusProgram() ([]byte, error) {
	return b.block.Body.NextConsensusProgram, nil
}

func (b *BlockVMContext) TxVersion() (uint64, bool)      { return 0, false }
func (b *BlockVMContext) TxSigHash() ([]byte, error)     { return nil, vm.ErrContext }
func (b *BlockVMContext) NumResults() (uint64, error)    { return 0, vm.ErrContext }
func (b *BlockVMContext) AssetID() ([]byte, error)       { return nil, vm.ErrContext }
func (b *BlockVMContext) Amount() (uint64, error)        { return 0, vm.ErrContext }
func (b *BlockVMContext) MinTimeMS() (uint64, error)     { return 0, vm.ErrContext }
func (b *BlockVMContext) MaxTimeMS() (uint64, error)     { return 0, vm.ErrContext }
func (b *BlockVMContext) EntryData() ([]byte, error)     { return nil, vm.ErrContext }
func (b *BlockVMContext) TxData() ([]byte, error)        { return nil, vm.ErrContext }
func (b *BlockVMContext) DestPos() (uint64, error)       { return 0, vm.ErrContext }
func (b *BlockVMContext) AnchorID() ([]byte, error)      { return nil, vm.ErrContext }
func (b *BlockVMContext) SpentOutputID() ([]byte, error) { return nil, vm.ErrContext }

func (b *BlockVMContext) CheckOutput(uint64, []byte, uint64, []byte, uint64, []byte) (bool, error) {
	return false, vm.ErrContext
}

func NewBlockVMContext(block *BlockEntries, prog []byte, args [][]byte) *BlockVMContext {
	return &BlockVMContext{
		prog: Program{
			VMVersion: 1,
			Code:      prog,
		},
		args:  args,
		block: block,
	}
}

type TxVMContext struct {
	prog  Program
	args  [][]byte
	tx    *TxEntries
	entry Entry
}

func NewTxVMContext(tx *TxEntries, entry Entry, prog Program, args [][]byte) *TxVMContext {
	return &TxVMContext{
		prog:  prog,
		args:  args,
		tx:    tx,
		entry: entry,
	}
}

func (t *TxVMContext) VMVersion() uint64   { return t.prog.VMVersion }
func (t *TxVMContext) Code() []byte        { return t.prog.Code }
func (t *TxVMContext) Arguments() [][]byte { return t.args }

func (t *TxVMContext) BlockHash() ([]byte, error)   { return nil, vm.ErrContext }
func (t *TxVMContext) BlockTimeMS() (uint64, error) { return 0, vm.ErrContext }

func (t *TxVMContext) NextConsensusProgram() ([]byte, error) { return nil, vm.ErrContext }

func (t *TxVMContext) TxVersion() (uint64, bool) { return t.tx.Body.Version, true }

func (t *TxVMContext) TxSigHash() ([]byte, error) {
	ord := t.entry.Ordinal()
	if ord < 0 {
		return nil, vm.ErrContext
	}
	h := t.tx.SigHash(uint32(ord))
	return h[:], nil
}

func (t *TxVMContext) NumResults() (uint64, error) { return uint64(len(t.tx.Results)), nil }

func (t *TxVMContext) AssetID() ([]byte, error) {
	switch inp := t.entry.(type) {
	case *Nonce:
		if iss, ok := inp.Anchored.(*Issuance); ok {
			return iss.Body.Value.AssetID[:], nil
		}
		return nil, vm.ErrContext

	case *Issuance:
		return inp.Body.Value.AssetID[:], nil

	case *Spend:
		return inp.SpentOutput.Body.Source.Value.AssetID[:], nil
	}

	return nil, vm.ErrContext
}

func (t *TxVMContext) Amount() (uint64, error) {
	switch inp := t.entry.(type) {
	case *Nonce:
		if iss, ok := inp.Anchored.(*Issuance); ok {
			return iss.Body.Value.Amount, nil
		}
		return 0, vm.ErrContext

	case *Issuance:
		return inp.Body.Value.Amount, nil

	case *Spend:
		return inp.SpentOutput.Body.Source.Value.Amount, nil
	}

	return 0, vm.ErrContext
}

func (t *TxVMContext) MinTimeMS() (uint64, error) { return t.tx.Body.MinTimeMS, nil }
func (t *TxVMContext) MaxTimeMS() (uint64, error) { return t.tx.Body.MaxTimeMS, nil }

func (t *TxVMContext) EntryData() ([]byte, error) {
	switch inp := t.entry.(type) {
	case *Issuance:
		return inp.Body.Data[:], nil

	case *Spend:
		return inp.Body.Data[:], nil

	case *Output:
		return inp.Body.Data[:], nil

	case *Retirement:
		return inp.Body.Data[:], nil
	}

	return nil, vm.ErrContext
}

func (t *TxVMContext) TxData() ([]byte, error) { return t.tx.Body.Data[:], nil }

func (t *TxVMContext) DestPos() (uint64, error) {
	switch inp := t.entry.(type) {
	case *Issuance:
		return inp.Witness.Destination.Position, nil

	case *Spend:
		return inp.Witness.Destination.Position, nil
	}

	return 0, vm.ErrContext
}

func (t *TxVMContext) AnchorID() ([]byte, error) {
	if inp, ok := t.entry.(*Issuance); ok {
		return inp.Body.AnchorID[:], nil
	}
	return nil, vm.ErrContext
}

func (t *TxVMContext) SpentOutputID() ([]byte, error) {
	if inp, ok := t.entry.(*Spend); ok {
		return inp.Body.SpentOutputID[:], nil
	}
	return nil, vm.ErrContext
}

func (t *TxVMContext) CheckOutput(index uint64, data []byte, amount uint64, assetID []byte, vmVersion uint64, code []byte) (bool, error) {
	checkEntry := func(e Entry) (bool, error) {
		check := func(prog Program, value AssetAmount, dataHash Hash) bool {
			return (prog.VMVersion == vmVersion &&
				bytes.Equal(prog.Code, code) &&
				bytes.Equal(value.AssetID[:], assetID) &&
				value.Amount == amount &&
				(len(data) == 0 || bytes.Equal(dataHash[:], data)))
		}

		switch e := e.(type) {
		case *Output:
			return check(e.Body.ControlProgram, e.Body.Source.Value, e.Body.Data), nil

		case *Retirement:
			return check(Program{}, e.Body.Source.Value, e.Body.Data), nil
		}

		return false, vm.ErrContext
	}

	checkMux := func(m *Mux) (bool, error) {
		if index >= uint64(len(m.Witness.Destinations)) {
			return errors.Wrapf(vm.ErrBadValue, "index %d >= %d", index, len(m.Witness.Destinations))
		}
		return checkEntry(m.Witness.Destinations[index].Entry)
	}

	switch e := t.entry.(type) {
	case *Mux:
		return checkMux(e)

	case *Issuance:
		if m, ok := e.Witness.Destination.Entry.(*Mux); ok {
			return checkMux(m)
		}
		if index != 0 {
			return errors.Wrapf(vm.ErrBadValue, "index %d >= 1", index)
		}
		return checkEntry(e.Witness.Destination.Entry)

	case *Spend:
		if m, ok := e.Witness.Destination.Entry.(*Mux); ok {
			return checkMux(m)
		}
		if index != 0 {
			return errors.Wrapf(vm.ErrBadValue, "index %d >= 1", index)
		}
		return checkEntry(e.Witness.Destination.Entry)
	}

	return false, vm.ErrContext
}
