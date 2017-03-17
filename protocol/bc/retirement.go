package bc

import (
	"chain/errors"
	"context"
)

// Retirement is for the permanent removal of some value from a
// blockchain. The value it contains can never be obtained by later
// entries. Retirement satisfies the Entry interface.
type Retirement struct {
	Body struct {
		Source  ValueSource
		Data    Hash
		ExtHash Hash
	}
	ordinal int
}

func (Retirement) Type() string         { return "retirement1" }
func (r *Retirement) body() interface{} { return r.Body }

func (r Retirement) Ordinal() int { return r.ordinal }

// NewRetirement creates a new Retirement.
func NewRetirement(source ValueSource, data Hash, ordinal int) *Retirement {
	r := new(Retirement)
	r.Body.Source = source
	r.Body.Data = data
	r.ordinal = ordinal
	return r
}

func (r *Retirement) CheckValid(ctx context.Context) error {
	ctx = context.WithValue(ctx, vcSourcePos, 0)
	err := r.Body.Source.CheckValid(ctx)
	if err != nil {
		return errors.Wrap(err, "checking retirement source")
	}

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	if currentTx.Body.Version == 1 && (r.Body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
