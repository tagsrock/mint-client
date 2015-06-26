package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
)

//------------------------------------------------------------------------------
// cli wrappers

func cliGenesis(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		Exit(fmt.Errorf("Enter the number of validators to create"))
	}
	N, err := strconv.Atoi(args[0])
	if err != nil {
		Exit(fmt.Errorf("Please provide an integer number of validators to create"))
	}

	fmt.Println("Generating accounts ...")
	genDoc, _, validators := state.RandGenesisState(N, true, 100000, N, false, 1000)

	name := NameFlag
	genDoc.ChainID = name

	genesisBytes, err := json.Marshal(genDoc)
	IfExit(err)

	// create directory to save priv validators and genesis.json
	d := path.Join(DataContainersPath, name)
	IfExit(os.MkdirAll(d, 0700))
	IfExit(ioutil.WriteFile(path.Join(d, "genesis.json"), genesisBytes, 0644))
	for i, v := range validators {
		valBytes, err := json.Marshal(v)
		IfExit(err)
		IfExit(ioutil.WriteFile(path.Join(d, fmt.Sprintf("priv_validator_%d.json", i)), valBytes, 0600))
	}
	fmt.Printf("genesis.json and priv_validator_X.json files saved in %s\n", d)
}
