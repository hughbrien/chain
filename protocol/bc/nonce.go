package bc

import (
	"chain/errors"
	"chain/protocol/vm"
	"context"
)

// Nonce contains data used, among other things, for distinguishing
// otherwise-identical issuances (when used as those issuances'
// "anchors"). It satisfies the Entry interface.
type Nonce struct {
	Body struct {
		Program     Program
		TimeRangeID Hash
		ExtHash     Hash
	}

	Witness struct {
		Arguments  [][]byte
		AnchoredID Hash
	}

	// TimeRange contains (a pointer to) the manifested entry
	// corresponding to body.TimeRangeID
	TimeRange *TimeRange

	// Anchored contains a pointer to the manifested entry corresponding
	// to witness.AnchoredID.
	Anchored Entry
}

func (Nonce) Type() string         { return "nonce1" }
func (n *Nonce) body() interface{} { return n.Body }

func (Nonce) Ordinal() int { return -1 }

// NewNonce creates a new Nonce.
func NewNonce(p Program, tr *TimeRange) *Nonce {
	n := new(Nonce)
	n.Body.Program = p
	n.Body.TimeRangeID = EntryID(tr)
	n.TimeRange = tr
	return n
}

func (n *Nonce) SetAnchored(id Hash, entry Entry) {
	n.Witness.AnchoredID = id
	n.Anchored = entry
}

func (n *Nonce) CheckValid(ctx context.Context) error {
	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	err := vm.Verify(NewTxVMContext(currentTx, n, n.Body.Program, n.Witness.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking nonce program")
	}

	trCtx := context.WithValue(ctx, vcCurrentEntryID, n.Body.TimeRangeID)
	err = n.TimeRange.CheckValid(trCtx)
	if err != nil {
		return errors.Wrap(err, "checking nonce timerange")
	}

	if n.TimeRange.Body.MinTimeMS == 0 || n.TimeRange.Body.MaxTimeMS == 0 {
		return errZeroTime
	}

	if currentTx.Body.Version == 1 && (n.Body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
