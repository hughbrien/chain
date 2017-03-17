package bc

import (
	"chain/errors"
	"chain/protocol/vm"
	"context"
)

type (
	NonceBody struct {
		Program     Program
		TimeRangeID Hash
		ExtHash     Hash
	}

	NonceWitness struct {
		Arguments  [][]byte
		AnchoredID Hash
	}

	// Nonce contains data used, among other things, for distinguishing
	// otherwise-identical issuances (when used as those issuances'
	// "anchors"). It satisfies the Entry interface.
	Nonce struct {
		NonceBody
		NonceWitness

		// TimeRange contains (a pointer to) the manifested entry
		// corresponding to body.TimeRangeID
		TimeRange *TimeRange

		// Anchored contains a pointer to the manifested entry corresponding
		// to witness.AnchoredID.
		Anchored Entry
	}
)

func (Nonce) Type() string         { return "nonce1" }
func (n *Nonce) Body() interface{} { return n.NonceBody }

func (Nonce) Ordinal() int { return -1 }

// NewNonce creates a new Nonce.
func NewNonce(p Program, tr *TimeRange) *Nonce {
	return &Nonce{
		NonceBody: NonceBody{
			Program:     p,
			TimeRangeID: EntryID(tr),
		},
	}
}

func (n *Nonce) CheckValid(ctx context.Context) error {
	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	err := vm.Verify(newTxVMContext(currentTx, n, n.Program, n.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking nonce program")
	}

	// xxx recursively validate the timerange?

	if n.TimeRange.MinTimeMS == 0 || n.TimeRange.MaxTimeMS == 0 {
		return errZeroTime
	}

	if currentTx.Version == 1 && (n.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
