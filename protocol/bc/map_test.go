package bc

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"chain/testutil"
)

func TestMapTx(t *testing.T) {
	// sample data copied from transaction_test.go
	// TODO(bobg): factor out into reusable test utility

	oldTx := sampleTx()
	oldOuts := oldTx.Outputs

	_, header, entryMap, err := mapTx(oldTx)
	if err != nil {
		testutil.FatalErr(t, err)
	}

	t.Log(spew.Sdump(entryMap))

	if header.Version != 1 {
		t.Errorf("header.Version is %d, expected 1", header.Version)
	}
	if header.MinTimeMS != oldTx.MinTime {
		t.Errorf("header.MinTimeMS is %d, expected %d", header.MinTimeMS, oldTx.MinTime)
	}
	if header.MaxTimeMS != oldTx.MaxTime {
		t.Errorf("header.MaxTimeMS is %d, expected %d", header.MaxTimeMS, oldTx.MaxTime)
	}
	if len(header.ResultIDs) != len(oldOuts) {
		t.Errorf("header.ResultIDs contains %d item(s), expected %d", len(header.ResultIDs), len(oldOuts))
	}

	for i, oldOut := range oldOuts {
		if resultEntry, ok := entryMap[header.ResultIDs[i]]; ok {
			if newOut, ok := resultEntry.(*Output); ok {
				if newOut.Source.Value != oldOut.AssetAmount {
					t.Errorf("header.ResultIDs[%d].(*output).Source is %v, expected %v", i, newOut.Source.Value, oldOut.AssetAmount)
				}
				if newOut.ControlProgram.VMVersion != 1 {
					t.Errorf("header.ResultIDs[%d].(*output).ControlProgram.VMVersion is %d, expected 1", i, newOut.ControlProgram.VMVersion)
				}
				if !bytes.Equal(newOut.ControlProgram.Code, oldOut.ControlProgram) {
					t.Errorf("header.ResultIDs[%d].(*output).ControlProgram.Code is %x, expected %x", i, newOut.ControlProgram.Code, oldOut.ControlProgram)
				}
				if newOut.Data != hashData(oldOut.ReferenceData) {
					want := hashData(oldOut.ReferenceData)
					t.Errorf("header.ResultIDs[%d].(*output).Data is %x, expected %x", i, newOut.Data[:], want[:])
				}
				if (newOut.ExtHash != Hash{}) {
					t.Errorf("header.ResultIDs[%d].(*output).ExtHash is %x, expected zero", i, newOut.ExtHash[:])
				}
			} else {
				t.Errorf("header.ResultIDs[%d] has type %s, expected output1", i, resultEntry.Type())
			}
		} else {
			t.Errorf("entryMap contains nothing for header.ResultIDs[%d] (%x)", i, header.ResultIDs[i][:])
		}
	}
}
