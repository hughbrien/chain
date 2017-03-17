package bc

import (
	"chain/errors"
	"chain/protocol/vm"
	"context"
)

type (
	SpendBody struct {
		SpentOutputID Hash // the hash of an output entry
		Data          Hash
		ExtHash       Hash
	}

	SpendWitness struct {
		Destination ValueDestination
		Arguments   [][]byte
		AnchoredID  Hash
	}

	// Spend accesses the value in a prior Output for transfer
	// elsewhere. It satisfies the Entry interface.
	//
	// (Not to be confused with the deprecated type SpendInput.)
	Spend struct {
		SpendBody
		SpendWitness

		ordinal int

		// SpentOutput contains (a pointer to) the manifested entry
		// corresponding to body.SpentOutputID.
		SpentOutput *Output

		// Anchored contains a pointer to the manifested entry corresponding
		// to witness.AnchoredID.
		Anchored Entry
	}
)

func (Spend) Type() string         { return "spend1" }
func (s *Spend) Body() interface{} { return s.SpendBody }

func (s Spend) Ordinal() int { return s.ordinal }

func (s *Spend) AssetID() AssetID {
	return s.SpentOutput.AssetID()
}

func (s *Spend) ControlProgram() Program {
	return s.SpentOutput.ControlProgram
}

func (s *Spend) Amount() uint64 {
	return s.SpentOutput.Amount()
}

func (s *Spend) SetDestination(id Hash, pos uint64, e Entry) {
	s.Destination = ValueDestination{
		Ref:      id,
		Position: pos,
		Entry:    e,
	}
}

// NewSpend creates a new Spend.
func NewSpend(out *Output, data Hash, ordinal int) *Spend {
	return &Spend{
		SpendBody: SpendBody{
			SpentOutputID: EntryID(out),
			Data:          data,
		},
		ordinal:     ordinal,
		SpentOutput: out,
	}
}

func (s *Spend) CheckValid(ctx context.Context) error {
	// xxx SpentOutput "present"

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	err := vm.Verify(newTxVMContext(currentTx, s, s.SpentOutput.ControlProgram, s.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking control program")
	}

	if s.SpentOutput.Source.Value != s.Destination.Value {
		return errors.WithDetailf(
			errMismatchedValue,
			"previous output is for %d unit(s) of %x, spend wants %d unit(s) of %x",
			s.SpentOutput.Source.Value.Amount,
			s.SpentOutput.Source.Value.AssetID[:],
			s.Destination.Value.Amount,
			s.Destination.Value.AssetID[:],
		)
	}

	ctx = context.WithValue(ctx, vcDestPos, 0)
	err = s.Destination.CheckValid(ctx)
	if err != nil {
		return errors.Wrap(err, "checking spend destination")
	}

	if currentTx.Version == 1 && (s.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
