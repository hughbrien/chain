package vm

type VMContext interface {
	VMVersion() uint64
	Code() []byte
	Arguments() [][]byte

	EntryID() []byte

	TxVersion() (uint64, bool)

	BlockHash() ([]byte, error)
	BlockTimeMS() (uint64, error)
	NextConsensusProgram() ([]byte, error)

	TxSigHash() ([]byte, error)
	NumResults() (uint64, error)
	AssetID() ([]byte, error)

	Amount() (uint64, error)
	MinTimeMS() (uint64, error)
	MaxTimeMS() (uint64, error)
	EntryData() ([]byte, error)
	TxData() ([]byte, error)
	DestPos() (uint64, error)
	AnchorID() ([]byte, error)
	SpentOutputID() ([]byte, error)
	CheckOutput(index uint64, data []byte, amount uint64, assetID []byte, vmVersion uint64, code []byte) (bool, error)
}