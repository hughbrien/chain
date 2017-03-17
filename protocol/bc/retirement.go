package bc

import (
	"chain/errors"
	"context"
)

type (
	RetirementBody struct {
		Source  ValueSource
		Data    Hash
		ExtHash Hash
	}

	// Retirement is for the permanent removal of some value from a
	// blockchain. The value it contains can never be obtained by later
	// entries. Retirement satisfies the Entry interface.
	Retirement struct {
		RetirementBody
		ordinal int
	}
)

func (Retirement) Type() string         { return "retirement1" }
func (r *Retirement) Body() interface{} { return r.RetirementBody }

func (r Retirement) Ordinal() int { return r.ordinal }

func (r *Retirement) AssetID() AssetID {
	return r.Source.Value.AssetID
}

func (r *Retirement) Amount() uint64 {
	return r.Source.Value.Amount
}

// NewRetirement creates a new Retirement.
func NewRetirement(source ValueSource, data Hash, ordinal int) *Retirement {
	return &Retirement{
		RetirementBody: RetirementBody{
			Source: source,
			Data:   data,
		},
		ordinal: ordinal,
	}
}

func (r *Retirement) CheckValid(ctx context.Context) error {
	ctx = context.WithValue(ctx, vcSourcePos, 0)
	err := r.Source.Entry.CheckValid(ctx)
	if err != nil {
		return errors.Wrap(err, "checking retirement source")
	}

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	if currentTx.Version == 1 && (r.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
