package prottest

import (
	"testing"
	"time"

	"chain/protocol/bc"
)

func TestNewIssuance(t *testing.T) {
	c := NewChain(t)
	iss := NewIssuanceTx(t, c)
	err := bc.ValidateTx(iss.TxEntries, 1, c.InitialBlockHash, bc.Millis(time.Now()))
	if err != nil {
		t.Error(err)
	}
}
