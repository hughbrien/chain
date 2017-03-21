package protocol_test

import (
	"context"
	"log"
	"testing"
	"time"

	. "chain/protocol"
	"chain/protocol/bc"
	"chain/protocol/prottest"
	"chain/testutil"
)

func TestRecoverSnapshotNoAdditionalBlocks(t *testing.T) {
	store := prottest.NewMemStore()
	b, err := NewInitialBlock(nil, 0, time.Now().Add(-time.Minute))
	if err != nil {
		testutil.FatalErr(t, err)
	}
	c1, err := NewChain(context.Background(), b.Hash(), store, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = c1.CommitAppliedBlock(context.Background(), b, NewSnapshot())
	if err != nil {
		testutil.FatalErr(t, err)
	}

	// Snapshots are applied asynchronously. This loops waits
	// until the snapshot is created.
	for {
		_, height, _ := store.LatestSnapshot(context.Background())
		if height > 0 {
			break
		}
	}

	ctx := context.Background()

	c2, err := NewChain(context.Background(), b.Hash(), store, nil)
	if err != nil {
		t.Fatal(err)
	}
	block, snapshot, err := c2.Recover(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if block.Height != 1 {
		t.Fatalf("block.Height = %d, want %d", block.Height, 1)
	}

	err = c2.ValidateBlockForSig(ctx, createEmptyBlock(block, snapshot))
	if err != nil {
		t.Fatal(err)
	}
}

func createEmptyBlock(block *bc.Block, snapshot *Snapshot) *bc.Block {
	root, err := bc.CalcMerkleRoot(nil)
	if err != nil {
		log.Fatalf("calculating empty merkle root: %s", err)
	}

	return &bc.Block{
		BlockHeader: bc.BlockHeader{
			Version:           bc.NewBlockVersion,
			Height:            block.Height + 1,
			PreviousBlockHash: block.Hash(),
			TimestampMS:       bc.Millis(time.Now()),
			BlockCommitment: bc.BlockCommitment{
				TransactionsMerkleRoot: root,
				AssetsMerkleRoot:       snapshot.Tree.RootHash(),
				ConsensusProgram:       block.ConsensusProgram,
			},
		},
	}
}
