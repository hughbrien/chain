package bc

import (
	"chain/errors"
	"chain/math/checked"
	"chain/protocol/vm"
	"context"
)

type (
	MuxBody struct {
		Sources []ValueSource // issuances, spends, and muxes
		Program Program
		ExtHash Hash
	}

	MuxWitness struct {
		Destinations []ValueDestination // outputs, retirements, and muxes
		Arguments    [][]byte
	}

	// Mux splits and combines value from one or more source entries,
	// making it available to one or more destination entries. It
	// satisfies the Entry interface.
	Mux struct {
		MuxBody
		MuxWitness
	}
)

func (Mux) Type() string         { return "mux1" }
func (m *Mux) Body() interface{} { return m.MuxBody }

func (Mux) Ordinal() int { return -1 }

// NewMux creates a new Mux.
func NewMux(sources []ValueSource, program Program) *Mux {
	return &Mux{
		MuxBody: MuxBody{
			Sources: sources,
			Program: program,
		},
	}
}

func (mux *Mux) CheckValid(ctx context.Context) error {
	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	err := vm.Verify(newTxVMContext(currentTx, mux, mux.Program, mux.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking mux program")
	}

	for i, src := range mux.Sources {
		ctx = context.WithValue(ctx, vcSourcePos, uint64(i))
		err := src.CheckValid(ctx)
		if err != nil {
			return errors.Wrapf(err, "checking mux source %d", i)
		}
	}

	for i, dest := range mux.Destinations {
		ctx = context.WithValue(ctx, vcDestPos, uint64(i))
		err := dest.CheckValid(ctx)
		if err != nil {
			return errors.Wrapf(err, "checking mux destination %d", i)
		}
	}

	parity := make(map[AssetID]int64)
	for i, src := range mux.Sources {
		sum, ok := checked.AddInt64(parity[src.Value.AssetID], int64(src.Value.Amount))
		if !ok {
			return errors.WithDetailf(errOverflow, "adding %d units of asset %x from mux source %d to total %d overflows int64", src.Value.Amount, src.Value.AssetID[:], i, parity[src.Value.AssetID])
		}
		parity[src.Value.AssetID] = sum
	}

	for i, dest := range mux.Destinations {
		sum, ok := parity[dest.Value.AssetID]
		if !ok {
			return errors.WithDetailf(errNoSource, "mux destination %d, asset %x, has no corresponding source", i, dest.Value.AssetID[:])
		}

		diff, ok := checked.SubInt64(sum, int64(dest.Value.Amount))
		if !ok {
			return errors.WithDetailf(errOverflow, "subtracting %d units of asset %x from mux destination %d from total %d underflows int64", dest.Value.Amount, dest.Value.AssetID[:], i, sum)
		}
		parity[dest.Value.AssetID] = diff
	}

	for assetID, amount := range parity {
		if amount != 0 {
			return errors.WithDetailf(errUnbalanced, "asset %x sources - destinations = %d (should be 0)", assetID[:], amount)
		}
	}

	if currentTx.Version == 1 && (mux.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
