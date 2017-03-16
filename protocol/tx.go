package protocol

import (
	"chain/errors"
	"chain/protocol/bc"
	"chain/protocol/validation"
)

func (c *Chain) checkIssuanceWindow(tx *bc.Tx) error {
	for _, txi := range tx.Inputs {
		if _, ok := txi.TypedInput.(*bc.IssuanceInput); ok {
			// TODO(tessr): consider removing 0 check once we can configure this
			if c.MaxIssuanceWindow != 0 && tx.MinTime+bc.DurationMillis(c.MaxIssuanceWindow) < tx.MaxTime {
				return errors.WithDetailf(validation.ErrBadTx, "issuance input's time window is larger than the network maximum (%s)", c.MaxIssuanceWindow)
			}
		}
	}
	return nil
}

func (c *Chain) ValidateTx(tx *bc.TxEntries, timestampMS uint64) error {
	return bc.ValidateTx(tx, bc.NewBlockVersion, c.InitialBlockHash, timestampMS)
}
