package bc

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMuxValid(t *testing.T) {
	var (
		mux *Mux
		vs  *validationState
	)

	cases := []struct {
		f   func()
		err error
	}{
		{},
		// xxx moar tests
	}

	for i, c := range cases {
		t.Logf("case %d", i)

		fixture := sample(t, nil)
		tx := NewTx(*fixture.tx)
		mux = tx.TxEntries.Results[0].(*Output).Body.Source.Entry.(*Mux)
		vs = &validationState{
			blockchainID: fixture.initialBlockID,
			tx:           tx.TxEntries,
			entryID:      tx.TxEntries.Results[0].(*Output).Body.Source.Ref,
		}

		if c.f != nil {
			c.f()
		}
		err := mux.CheckValid(vs)
		if err != c.err {
			t.Errorf("case %d: got error %s, want %s; mux is:\n%s\nvalidationState is:\n%s", i, err, c.err, spew.Sdump(mux), spew.Sdump(vs))
		}
	}
}
