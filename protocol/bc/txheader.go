package bc

import (
	"chain/errors"
	"context"
)

type (
	TxHeaderBody struct {
		Version              uint64
		ResultIDs            []Hash
		Data                 Hash
		MinTimeMS, MaxTimeMS uint64
		ExtHash              Hash
	}

	// TxHeader contains header information for a transaction. Every
	// transaction on a blockchain contains exactly one TxHeader. The ID
	// of the TxHeader is the ID of the transaction. TxHeader satisfies
	// the Entry interface.
	TxHeader struct {
		TxHeaderBody

		// Results contains (pointers to) the manifested entries for the
		// items in body.ResultIDs.
		Results []Entry // each entry is *output or *retirement
	}
)

func (TxHeader) Type() string         { return "txheader" }
func (h *TxHeader) Body() interface{} { return h.TxHeaderBody }

func (TxHeader) Ordinal() int { return -1 }

// NewTxHeader creates an new TxHeader.
func NewTxHeader(version uint64, results []Entry, data Hash, minTimeMS, maxTimeMS uint64) *TxHeader {
	h := &TxHeader{
		TxHeaderBody: TxHeaderBody{
			Version:   version,
			Data:      data,
			MinTimeMS: minTimeMS,
			MaxTimeMS: maxTimeMS,
		},
		Results: results,
	}
	for _, r := range results {
		h.ResultIDs = append(h.ResultIDs, EntryID(r))
	}
	return h
}

// CheckValid does only part of the work of validating a tx header. The block-related parts of tx validation are in ValidateBlock.
func (tx *TxHeader) CheckValid(ctx context.Context) error {
	if tx.MaxTimeMS > 0 {
		if tx.MaxTimeMS < tx.MinTimeMS {
			return errors.WithDetailf(errBadTimeRange, "min time %d, max time %d", tx.MinTimeMS, tx.MaxTimeMS)
		}
	}

	for i, resID := range tx.ResultIDs {
		res := tx.Results[i]
		ctx = context.WithValue(ctx, vcCurrentEntryID, resID)
		err := res.CheckValid(ctx)
		if err != nil {
			return errors.Wrapf(err, "checking result %d", i)
		}
	}

	if tx.Version == 1 {
		if len(tx.ResultIDs) == 0 {
			return errEmptyResults
		}

		if (tx.ExtHash != Hash{}) {
			return errNonemptyExtHash
		}
	}

	return nil
}
