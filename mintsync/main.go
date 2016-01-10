package main

import (
	bc "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/blockchain"
	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/common"
	dbm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/db"
	sm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
)

func main() {
	// Get BlockStore
	blockStoreDB := dbm.GetDB("blockstore")
	blockStore := bc.NewBlockStore(blockStoreDB)

	// Get State
	stateDB := dbm.GetDB("state")
	state := sm.LoadState(stateDB)

	// replay blocks on the state
	var block, nextBlock *types.Block
	if state.LastBlockHeight < blockStore.Height()-1 {
		for i := 1; i < blockStore.Height()-state.LastBlockHeight; i++ {
			block = blockStore.LoadBlock(state.LastBlockHeight + i)
			nextBlock = blockStore.LoadBlock(state.LastBlockHeight + i + 1)
			parts := block.MakePartSet()
			err := sm.ExecBlock(state, block, parts.Header())
			if err != nil {
				// TODO This is bad, are we zombie?
				PanicQ(Fmt("Failed to process committed block: %v", err))
			}
			state.Save()
			block = nextBlock
		}
	}
}
