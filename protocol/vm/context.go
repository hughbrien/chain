package vm

// Context contains the execution context for the virtual machine.
//
// By convention, variables of this type have the name context, _not_
// ctx (to avoid confusion with context.Context).
type Context struct {
	VMVersion uint64
	Code      []byte
	Arguments [][]byte

	EntryID []byte

	TxVersion *uint64

	BlockHash            *[]byte
	BlockTimeMS          *uint64
	NextConsensusProgram *[]byte

	NumResults    *uint64
	AssetID       *[]byte
	Amount        *uint64
	MinTimeMS     *uint64
	MaxTimeMS     *uint64
	EntryData     *[]byte
	TxData        *[]byte
	DestPos       *uint64
	AnchorID      *[]byte
	SpentOutputID *[]byte

	TxSigHash   func() []byte
	CheckOutput func(index uint64, data []byte, amount uint64, assetID []byte, vmVersion uint64, code []byte) (bool, error)
}
