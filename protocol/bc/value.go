package bc

import (
	"context"
	"fmt"

	"chain/errors"
)

type ValueSource struct {
	Ref      Hash
	Value    AssetAmount
	Position uint64

	// The Entry corresponding to Ref, if available
	// The struct tag excludes the field from hashing
	Entry `entry:"-"`
}

// CheckValid checks the validity of a value source in the context of
// its containing entry.
func (vs *ValueSource) CheckValid(ctx context.Context) error {
	// xxx check that Entry's ID equals Ref?

	ctx = context.WithValue(ctx, vcCurrentEntryID, vs.Ref)
	err := vs.Entry.CheckValid(ctx)
	if err != nil {
		return errors.Wrap(err, "checking value source")
	}

	var dest ValueDestination
	switch ref := vs.Entry.(type) {
	case *Issuance:
		if vs.Position != 0 {
			return errors.WithDetailf(errPosition, "invalid position %d for issuance source", vs.Position)
		}
		dest = ref.Destination

	case *Spend:
		if vs.Position != 0 {
			return errors.WithDetailf(errPosition, "invalid position %d for spend source", vs.Position)
		}
		dest = ref.Destination

	case *Mux:
		if vs.Position >= uint64(len(ref.Destinations)) {
			return errors.WithDetailf(errPosition, "invalid position %d for %d-destination mux source", vs.Position, len(ref.Destinations))
		}
		dest = ref.Destinations[vs.Position]

	default:
		return errors.WithDetailf(errEntryType, "value source is %T, should be issuance, spend, or mux", vs.Entry)
	}

	currentEntryID, _ := ctx.Value(vcCurrentEntryID).(Hash)
	if dest.Ref != currentEntryID {
		return errors.WithDetailf(errMismatchedReference, "value source for %x has disagreeing destination %x", currentEntryID[:], dest.Ref[:])
	}

	sourcePos, _ := ctx.Value(vcSourcePos).(uint64)
	if dest.Position != sourcePos {
		return fmt.Errorf("value source position %d disagrees with %d", dest.Position, sourcePos)
	}

	if dest.Value != vs.Value {
		return fmt.Errorf("source value %v disagrees with %v", dest.Value, vs.Value)
	}

	return nil
}

type ValueDestination struct {
	Ref      Hash
	Value    AssetAmount
	Position uint64

	// The Entry corresponding to Ref, if available
	// The struct tag excludes the field from hashing
	Entry `entry:"-"`
}

func (vd *ValueDestination) CheckValid(ctx context.Context) error {
	// xxx check reachability of vd.Ref from transaction?

	var src ValueSource
	switch ref := vd.Entry.(type) {
	case *Output:
		if vd.Position != 0 {
			fmt.Errorf("invalid position %d for output destination", vd.Position)
		}
		src = ref.Source

	case *Retirement:
		if vd.Position != 0 {
			fmt.Errorf("invalid position %d for retirement destination", vd.Position)
		}
		src = ref.Source

	case *Mux:
		if vd.Position >= uint64(len(ref.Sources)) {
			return fmt.Errorf("invalid position %d for %d-source mux destination", vd.Position, len(ref.Sources))
		}
		src = ref.Sources[vd.Position]

	default:
		return fmt.Errorf("value destination is %T, should be output, retirement, or mux", vd.Entry)
	}

	currentEntryID, _ := ctx.Value(vcCurrentEntryID).(Hash)
	if src.Ref != currentEntryID {
		return fmt.Errorf("value destination for %x has disagreeing source %x", currentEntryID[:], src.Ref[:])
	}

	destPos, _ := ctx.Value(vcDestPos).(uint64)
	if src.Position != destPos {
		return fmt.Errorf("value destination position %d disagrees with %d", src.Position, destPos)
	}

	if src.Value != vd.Value {
		return fmt.Errorf("destination value %v disagrees with %v", src.Value, vd.Value)
	}

	return nil
}
