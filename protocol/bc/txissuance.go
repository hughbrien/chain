package bc

import (
	"chain/errors"
	"chain/protocol/vm"
	"context"
)

type (
	IssuanceBody struct {
		AnchorID Hash
		Value    AssetAmount
		Data     Hash
		ExtHash  Hash
	}

	IssuanceEntryWitness struct { // TODO(bobg): rename to IssuanceWitness when it no longer conflicts with the legacy type
		Destination     ValueDestination
		AssetDefinition AssetDefinition
		Arguments       [][]byte
		AnchoredID      Hash
	}

	// Issuance is a source of new value on a blockchain. It satisfies the
	// Entry interface.
	//
	// (Not to be confused with the deprecated type IssuanceInput.)
	Issuance struct {
		IssuanceBody
		IssuanceEntryWitness

		ordinal int

		// Anchor is a pointer to the manifested entry corresponding to
		// body.AnchorID.
		Anchor Entry // *nonce, *spend, or *issuance

		// Anchored is a pointer to the manifested entry corresponding to
		// witness.AnchoredID.
		Anchored Entry
	}
)

func (Issuance) Type() string           { return "issuance1" }
func (iss *Issuance) Body() interface{} { return iss.IssuanceBody }

func (iss Issuance) Ordinal() int { return iss.ordinal }

func (iss *Issuance) SetDestination(id Hash, pos uint64, e Entry) {
	iss.Destination = ValueDestination{
		Ref:      id,
		Position: pos,
		Entry:    e,
	}
}

func (iss *Issuance) SetInitialBlockID(hash Hash) {
	iss.AssetDefinition.InitialBlockID = hash
}

func (iss *Issuance) SetAssetDefinitionHash(hash Hash) {
	iss.AssetDefinition.Data = hash
}

func (iss *Issuance) SetIssuanceProgram(prog Program) {
	iss.AssetDefinition.IssuanceProgram = prog
}

// NewIssuance creates a new Issuance.
func NewIssuance(anchor Entry, value AssetAmount, data Hash, ordinal int) *Issuance {
	return &Issuance{
		IssuanceBody: IssuanceBody{
			AnchorID: EntryID(anchor),
			Value:    value,
			Data:     data,
		},
		Anchor:  anchor,
		ordinal: ordinal,
	}
}

func (iss *Issuance) CheckValid(ctx context.Context) error {
	initialBlockID, _ := ctx.Value(vcInitialBlockID).(Hash)
	if iss.AssetDefinition.InitialBlockID != initialBlockID {
		return errors.WithDetailf(errWrongBlockchain, "current blockchain %x, asset defined on blockchain %x", initialBlockID[:], iss.AssetDefinition.InitialBlockID[:])
	}

	computedAssetID := iss.AssetDefinition.ComputeAssetID()
	if computedAssetID != iss.Value.AssetID {
		return errors.WithDetailf(errMismatchedAssetID, "asset ID is %x, issuance wants %x", computedAssetID[:], iss.Value.AssetID[:])
	}

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	err := vm.Verify(newTxVMContext(currentTx, iss, iss.AssetDefinition.IssuanceProgram, iss.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking issuance program")
	}

	var anchored Hash
	switch a := iss.Anchor.(type) {
	case *Nonce:
		anchored = a.AnchoredID

	case *Spend:
		anchored = a.AnchoredID

	case *Issuance:
		anchored = a.AnchoredID

	default:
		return errors.WithDetailf(errEntryType, "issuance anchor has type %T, should be nonce, spend, or issuance", iss.Anchor)
	}

	currentEntryID, _ := ctx.Value(vcCurrentEntryID).(Hash)
	if anchored != currentEntryID {
		return errors.WithDetailf(errMismatchedReference, "issuance %x anchor is for %x", currentEntryID[:], anchored[:])
	}

	anchorCtx := context.WithValue(ctx, vcCurrentEntryID, iss.AnchorID)
	err = iss.Anchor.CheckValid(anchorCtx)
	if err != nil {
		return errors.Wrap(err, "checking issuance anchor")
	}

	destCtx := context.WithValue(ctx, vcDestPos, 0)
	err = iss.Destination.CheckValid(destCtx)
	if err != nil {
		return errors.Wrap(err, "checking issuance destination")
	}

	if currentTx.Version == 1 && (iss.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
