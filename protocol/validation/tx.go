package validation

import (
	"chain/errors"
	"chain/protocol"
	"chain/protocol/bc"
)

// ErrBadTx is returned for transactions failing validation
var ErrBadTx = errors.New("invalid transaction")

var (
	// "suberrors" for ErrBadTx
	errTxVersion              = errors.New("unknown transaction version")
	errNotYet                 = errors.New("block time is before transaction min time")
	errTooLate                = errors.New("block time is after transaction max time")
	errWrongBlockchain        = errors.New("issuance is for different blockchain")
	errTimelessIssuance       = errors.New("zero mintime or maxtime not allowed in issuance with non-empty nonce")
	errIssuanceTime           = errors.New("timestamp outside issuance input's time window")
	errDuplicateNonce         = errors.New("duplicate nonce entry")
	errInvalidOutput          = errors.New("invalid output")
	errNoInputs               = errors.New("inputs are missing")
	errTooManyInputs          = errors.New("number of inputs overflows uint32")
	errAllEmptyNonceIssuances = errors.New("all inputs are issuances with empty nonce fields")
	errMisorderedTime         = errors.New("positive maxtime must be >= mintime")
	errAssetVersion           = errors.New("unknown asset version")
	errInputTooBig            = errors.New("input value exceeds maximum value of int64")
	errInputSumTooBig         = errors.New("sum of inputs overflows the allowed asset amount")
	errVMVersion              = errors.New("unknown vm version")
	errDuplicateInput         = errors.New("duplicate input")
	errTooManyOutputs         = errors.New("number of outputs overflows int32")
	errEmptyOutput            = errors.New("output value must be greater than 0")
	errOutputTooBig           = errors.New("output value exceeds maximum value of int64")
	errOutputSumTooBig        = errors.New("sum of outputs overflows the allowed asset amount")
	errUnbalancedV1           = errors.New("amounts for asset are not balanced on v1 inputs and outputs")
)

func badTxErr(err error) error {
	err = errors.WithData(err, "badtx", err)
	err = errors.WithDetail(err, err.Error())
	return errors.Sub(ErrBadTx, err)
}

func badTxErrf(err error, f string, args ...interface{}) error {
	err = errors.WithData(err, "badtx", err)
	err = errors.WithDetailf(err, f, args...)
	return errors.Sub(ErrBadTx, err)
}

// ConfirmTx validates the given transaction against the given state tree
// before it's added to a block. If tx is invalid, it returns a non-nil
// error describing why.
//
// Tx must already have undergone the well-formedness check in
// CheckTxWellFormed. This should have happened when the tx was added
// to the pool.
//
// ConfirmTx must not mutate the snapshot.
func ConfirmTx(snapshot *protocol.Snapshot, initialBlockHash bc.Hash, blockVersion, blockTimestampMS uint64, tx *bc.TxEntries) error {
	if tx.Version() < 1 || tx.Version() > blockVersion {
		return badTxErrf(errTxVersion, "unknown transaction version %d for block version %d", tx.Version, blockVersion)
	}

	if blockTimestampMS < tx.MinTimeMS() {
		return badTxErr(errNotYet)
	}
	if tx.MaxTimeMS() > 0 && blockTimestampMS > tx.MaxTimeMS() {
		return badTxErr(errTooLate)
	}

	for i, inp := range tx.TxInputs {
		switch inp := inp.(type) {
		case *bc.Issuance:
			if inp.InitialBlockID() != initialBlockHash {
				return badTxErr(errWrongBlockchain)
			}
			// xxx nonce/timerange check (already done in checktxwellformed)?
			if blockTimestampMS < tx.MinTimeMS() || blockTimestampMS > tx.MaxTimeMS() {
				return badTxErr(errIssuanceTime)
			}
			id := tx.TxInputIDs[i]
			if _, ok := snapshot.Nonces[id]; ok {
				return badTxErr(errDuplicateNonce)
			}

		case *bc.Spend:
			if !snapshot.Tree.Contains(inp.SpentOutputID().Bytes()) {
				return badTxErrf(errInvalidOutput, "output %s for input %d is invalid", inp.SpentOutputID(), i)
			}
		}
	}
	return nil
}

// ApplyTx updates the state tree with all the changes to the ledger.
func ApplyTx(snapshot *protocol.Snapshot, tx *bc.TxEntries) error {
	for i, inp := range tx.TxInputs {
		switch inp := inp.(type) {
		case *bc.Issuance:
			id := tx.TxInputIDs[i]
			snapshot.Issuances[id] = tx.MaxTimeMS() // xxx or the max time from the anchor timerange?

		case *bc.Spend:
			// Remove the consumed output from the state tree.
			snapshot.Tree.Delete(inp.SpentOutputID().Bytes())
		}
	}

	for i, res := range tx.Results {
		if _, ok := res.(*bc.Output); ok {
			err := snapshot.Tree.Insert(tx.ResultID(uint32(i)).Bytes())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
