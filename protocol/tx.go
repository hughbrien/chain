package protocol

import (
	"chain/errors"
	"chain/protocol/bc"
)

// ErrBadTx is returned for transactions failing validation
var ErrBadTx = errors.New("invalid transaction")

func (c *Chain) checkIssuanceWindow(tx *bc.TxEntries) error {
	if c.MaxIssuanceWindow == 0 {
		return nil
	}
	for _, entry := range tx.TxInputs {
		if _, ok := entry.(*bc.Issuance); ok {
			if tx.Body.MinTimeMS+bc.DurationMillis(c.MaxIssuanceWindow) < tx.Body.MaxTimeMS {
				return errors.WithDetailf(ErrBadTx, "issuance input's time window is larger than the network maximum (%s)", c.MaxIssuanceWindow)
			}
		}
	}
	return nil
}

func (c *Chain) ValidateTx(tx *bc.TxEntries) error {
	err := c.checkIssuanceWindow(tx)
	if err != nil {
		return err
	}
	return bc.ValidateTx(tx, c.InitialBlockHash)
}
