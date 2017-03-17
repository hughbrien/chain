package bc

import (
	"chain/errors"
	"context"
)

type (
	OutputBody struct {
		Source         ValueSource
		ControlProgram Program
		Data           Hash
		ExtHash        Hash
	}

	// Output is the result of a transfer of value. The value it contains
	// may be accessed by a later Spend entry (if that entry can satisfy
	// the Output's ControlProgram). Output satisfies the Entry interface.
	//
	// (Not to be confused with the deprecated type TxOutput.)
	Output struct {
		OutputBody
		ordinal int
	}
)

func (Output) Type() string         { return "output1" }
func (o *Output) Body() interface{} { return o.OutputBody }

func (o Output) Ordinal() int { return o.ordinal }

func (o *Output) AssetID() AssetID {
	return o.Source.Value.AssetID
}

func (o *Output) Amount() uint64 {
	return o.Source.Value.Amount
}

func (o *Output) SourceID() Hash {
	return o.Source.Ref
}

func (o *Output) SourcePosition() uint64 {
	return o.Source.Position
}

// NewOutput creates a new Output.
func NewOutput(source ValueSource, controlProgram Program, data Hash, ordinal int) *Output {
	return &Output{
		OutputBody: OutputBody{
			Source:         source,
			ControlProgram: controlProgram,
			Data:           data,
		},
		ordinal: ordinal,
	}
}

func (o *Output) CheckValid(ctx context.Context) error {
	ctx = context.WithValue(ctx, vcSourcePos, 0)
	err := o.Source.CheckValid(ctx)
	if err != nil {
		return errors.Wrap(err, "checking output source")
	}

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	if currentTx.Version == 1 && (o.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
