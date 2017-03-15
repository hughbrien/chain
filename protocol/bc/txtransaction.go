package bc

import (
	"fmt"

	"chain/crypto/sha3pool"
	"chain/errors"
)

type TxEntries struct {
	*TxHeader
	ID         Hash
	TxInputs   []Entry // 1:1 correspondence with TxData.Inputs
	TxInputIDs []Hash  // 1:1 correspondence with TxData.Inputs
}

// ValidateTx validates a transaction; so does
// TxEntries.CheckValid. This one is more suitable for calling from
// the top level; CheckValid is preferred in a nested validation
// context (such as when validating all the transactions in a block).
func ValidateTx(tx *TxEntries, blockVersion uint64, initialBlockID Hash, blockTimestampMS uint64) error {
	state := &validationState{
		blockVersion: blockVersion,
		initialBlockID: initialBlockID,
		currentTx: tx,
		currentEntryID: tx.ID,
		timestampMS: blockTimestampMS,
	}
	return tx.TxHeader.CheckValid(state)
}

func (tx *TxEntries) CheckValid(state *validationState) error {
	var newState validationState
	if state != nil {
		newState = *state
	}
	newState.currentTx = tx
	newState.currentEntryID = tx.ID
	return tx.TxHeader.CheckValid(&newState)
}

func (tx *TxEntries) SigHash(n uint32) (hash Hash) {
	hasher := sha3pool.Get256()
	defer sha3pool.Put256(hasher)

	hasher.Write(tx.TxInputIDs[n][:])
	hasher.Write(tx.ID[:])

	hasher.Read(hash[:])
	return hash
}

// ComputeOutputID assembles an output entry given a spend commitment
// and computes and returns its corresponding entry ID.
func ComputeOutputID(sc *SpendCommitment) (h Hash, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()
	src := valueSource{
		Ref:      sc.SourceID,
		Value:    sc.AssetAmount,
		Position: sc.SourcePosition,
	}
	o := NewOutput(src, Program{VMVersion: sc.VMVersion, Code: sc.ControlProgram}, sc.RefDataHash, 0)

	h = EntryID(o)
	return h, nil
}

// TxHashes returns all hashes needed for validation and state updates.
func ComputeTxEntries(oldTx *TxData) (txEntries *TxEntries, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()

	txid, header, entries, err := mapTx(oldTx)
	if err != nil {
		return nil, errors.Wrap(err, "mapping old transaction to new")
	}

	txEntries = &TxEntries{
		TxHeader:   header,
		ID:         txid,
		TxInputs:   make([]Entry, len(oldTx.Inputs)),
		TxInputIDs: make([]Hash, len(oldTx.Inputs)),
	}

	for id, e := range entries {
		switch e.(type) {
		case *Issuance:
		case *Spend:
		default:
			continue
		}
		ord := e.Ordinal()
		if ord < 0 || ord >= len(oldTx.Inputs) {
			return nil, fmt.Errorf("%T entry has out-of-range ordinal %d", e, ord)
		}
		txEntries.TxInputs[ord] = e
		txEntries.TxInputIDs[ord] = id
	}

	return txEntries, nil
}
