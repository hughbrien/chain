package bc

import (
	"chain/errors"
	"context"
)

// TxHeader contains header information for a transaction. Every
// transaction on a blockchain contains exactly one TxHeader. The ID
// of the TxHeader is the ID of the transaction. TxHeader satisfies
// the Entry interface.
type TxHeader struct {
	body struct {
		Version              uint64
		ResultIDs            []Hash
		Data                 Hash
		MinTimeMS, MaxTimeMS uint64
		ExtHash              Hash
	}

	// Results contains (pointers to) the manifested entries for the
	// items in body.ResultIDs.
	Results []Entry // each entry is *output or *retirement
}

func (TxHeader) Type() string         { return "txheader" }
func (h *TxHeader) Body() interface{} { return h.body }

func (TxHeader) Ordinal() int { return -1 }

func (h *TxHeader) Version() uint64 {
	return h.body.Version
}

func (h *TxHeader) Data() Hash {
	return h.body.Data
}

func (h *TxHeader) ResultID(n uint32) Hash {
	return h.body.ResultIDs[n]
}

func (h *TxHeader) MinTimeMS() uint64 {
	return h.body.MinTimeMS
}

func (h *TxHeader) MaxTimeMS() uint64 {
	return h.body.MaxTimeMS
}

// NewTxHeader creates an new TxHeader.
func NewTxHeader(version uint64, results []Entry, data Hash, minTimeMS, maxTimeMS uint64) *TxHeader {
	h := new(TxHeader)
	h.body.Version = version
	h.body.Data = data
	h.body.MinTimeMS = minTimeMS
	h.body.MaxTimeMS = maxTimeMS

	h.Results = results
	for _, r := range results {
		h.body.ResultIDs = append(h.body.ResultIDs, EntryID(r))
	}

	return h
}

// CheckValid does only part of the work of validating a tx header. The block-related parts of tx validation are in ValidateBlock.
func (tx *TxHeader) CheckValid(ctx context.Context) error {
	if tx.body.MaxTimeMS > 0 {
		if tx.body.MaxTimeMS < tx.body.MinTimeMS {
			return errors.WithDetailf(errBadTimeRange, "min time %d, max time %d", tx.body.MinTimeMS, tx.body.MaxTimeMS)
		}
	}

	for i, resID := range tx.body.ResultIDs {
		res := tx.Results[i]
		ctx = context.WithValue(ctx, vcCurrentEntryID, resID)
		err := res.CheckValid(ctx)
		if err != nil {
			return errors.Wrapf(err, "checking result %d", i)
		}
	}

	if tx.body.Version == 1 {
		if len(tx.body.ResultIDs) == 0 {
			return errEmptyResults
		}

		if (tx.body.ExtHash != Hash{}) {
			return errNonemptyExtHash
		}
	}

	return nil
}
