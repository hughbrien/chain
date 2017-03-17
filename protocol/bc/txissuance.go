package bc

import (
	"chain/errors"
	"chain/protocol/vm"
	"context"
)

// Issuance is a source of new value on a blockchain. It satisfies the
// Entry interface.
//
// (Not to be confused with the deprecated type IssuanceInput.)
type Issuance struct {
	body struct {
		Anchor  Hash
		Value   AssetAmount
		Data    Hash
		ExtHash Hash
	}
	ordinal int

	witness struct {
		Destination     ValueDestination
		AssetDefinition AssetDefinition
		Arguments       [][]byte
		Anchored        Hash
	}

	// Anchor is a pointer to the manifested entry corresponding to
	// body.Anchor.
	Anchor Entry // *nonce, *spend, or *issuance

	// Anchored is a pointer to the manifested entry corresponding to
	// witness.Anchored.
	Anchored Entry
}

func (Issuance) Type() string           { return "issuance1" }
func (iss *Issuance) Body() interface{} { return iss.body }

func (iss Issuance) Ordinal() int { return iss.ordinal }

func (iss *Issuance) AnchorID() Hash {
	return iss.body.Anchor
}

func (iss *Issuance) Data() Hash {
	return iss.body.Data
}

func (iss *Issuance) AssetID() AssetID {
	return iss.body.Value.AssetID
}

func (iss *Issuance) Amount() uint64 {
	return iss.body.Value.Amount
}

func (iss *Issuance) Destination() ValueDestination {
	return iss.witness.Destination
}

func (iss *Issuance) InitialBlockID() Hash {
	return iss.witness.AssetDefinition.InitialBlockID
}

func (iss *Issuance) IssuanceProgram() Program {
	return iss.witness.AssetDefinition.IssuanceProgram
}

func (iss *Issuance) Arguments() [][]byte {
	return iss.witness.Arguments
}

func (iss *Issuance) SetDestination(id Hash, pos uint64, e Entry) {
	iss.witness.Destination = ValueDestination{
		Ref:      id,
		Position: pos,
		Entry:    e,
	}
}

func (iss *Issuance) SetInitialBlockID(hash Hash) {
	iss.witness.AssetDefinition.InitialBlockID = hash
}

func (iss *Issuance) SetAssetDefinitionHash(hash Hash) {
	iss.witness.AssetDefinition.Data = hash
}

func (iss *Issuance) SetIssuanceProgram(prog Program) {
	iss.witness.AssetDefinition.IssuanceProgram = prog
}

func (iss *Issuance) SetArguments(args [][]byte) {
	iss.witness.Arguments = args
}

// NewIssuance creates a new Issuance.
func NewIssuance(anchor Entry, value AssetAmount, data Hash, ordinal int) *Issuance {
	iss := new(Issuance)
	iss.body.Anchor = EntryID(anchor)
	iss.Anchor = anchor
	iss.body.Value = value
	iss.body.Data = data
	iss.ordinal = ordinal
	return iss
}

func (iss *Issuance) CheckValid(ctx context.Context) error {
	initialBlockID, _ := ctx.Value(vcInitialBlockID).(Hash)
	if iss.witness.AssetDefinition.InitialBlockID != initialBlockID {
		return errors.WithDetailf(errWrongBlockchain, "current blockchain %x, asset defined on blockchain %x", initialBlockID[:], iss.witness.AssetDefinition.InitialBlockID[:])
	}

	computedAssetID := iss.witness.AssetDefinition.ComputeAssetID()
	if computedAssetID != iss.body.Value.AssetID {
		return errors.WithDetailf(errMismatchedAssetID, "asset ID is %x, issuance wants %x", computedAssetID[:], iss.body.Value.AssetID[:])
	}

	currentTx, _ := ctx.Value(vcCurrentTx).(*TxEntries)
	err := vm.Verify(newTxVMContext(currentTx, iss, iss.witness.AssetDefinition.IssuanceProgram, iss.witness.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking issuance program")
	}

	var anchored Hash
	switch a := iss.Anchor.(type) {
	case *Nonce:
		anchored = a.witness.Anchored

	case *Spend:
		anchored = a.witness.Anchored

	case *Issuance:
		anchored = a.witness.Anchored

	default:
		return errors.WithDetailf(errEntryType, "issuance anchor has type %T, should be nonce, spend, or issuance", iss.Anchor)
	}

	currentEntryID, _ := ctx.Value(vcCurrentEntryID).(Hash)
	if anchored != currentEntryID {
		return errors.WithDetailf(errMismatchedReference, "issuance %x anchor is for %x", currentEntryID[:], anchored[:])
	}

	anchorCtx := context.WithValue(ctx, vcCurrentEntryID, iss.body.Anchor)
	err = iss.Anchor.CheckValid(anchorCtx)
	if err != nil {
		return errors.Wrap(err, "checking issuance anchor")
	}

	destCtx := context.WithValue(ctx, vcDestPos, 0)
	err = iss.witness.Destination.CheckValid(destCtx)
	if err != nil {
		return errors.Wrap(err, "checking issuance destination")
	}

	if currentTx.body.Version == 1 && (iss.body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
