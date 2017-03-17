package bc

import (
	"chain/errors"
	"context"
)

// Output is the result of a transfer of value. The value it contains
// may be accessed by a later Spend entry (if that entry can satisfy
// the Output's ControlProgram). Output satisfies the Entry interface.
//
// (Not to be confused with the deprecated type TxOutput.)
type Output struct {
	body struct {
		Source         valueSource
		ControlProgram Program
		Data           Hash
		ExtHash        Hash
	}
	ordinal int
}

func (Output) Type() string         { return "output1" }
func (o *Output) Body() interface{} { return o.body }

func (o Output) Ordinal() int { return o.ordinal }

func (o *Output) AssetID() AssetID {
	return o.body.Source.Value.AssetID
}

func (o *Output) Amount() uint64 {
	return o.body.Source.Value.Amount
}

func (o *Output) SourceID() Hash {
	return o.body.Source.Ref
}

func (o *Output) SourcePosition() uint64 {
	return o.body.Source.Position
}

func (o *Output) ControlProgram() Program {
	return o.body.ControlProgram
}

func (o *Output) Data() Hash {
	return o.body.Data
}

// NewOutput creates a new Output.
func NewOutput(source valueSource, controlProgram Program, data Hash, ordinal int) *Output {
	out := new(Output)
	out.body.Source = source
	out.body.ControlProgram = controlProgram
	out.body.Data = data
	out.ordinal = ordinal
	return out
}

func (o *Output) CheckValid(ctx context.Context) error {
	ctx = context.WithValue(ctx, vcSourcePos, 0)
	err := o.body.Source.CheckValid(ctx)
	if err != nil {
		return errors.Wrap(err, "checking output source")
	}

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	if currentTx.body.Version == 1 && (o.body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
