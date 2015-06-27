package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/binary"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
)

//------------------------------------------------------------------------------

var defaultConfig = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

moniker = "__MONIKER__"
node_laddr = "0.0.0.0:46656"
seeds = ""
fast_sync = false
db_backend = "leveldb"
log_level = "debug"
rpc_laddr = "0.0.0.0:46657"
`

func singleUserChain(dir string) {
	genDoc, _, validators := state.RandGenesisDoc(1, true, 1000000000, 1, false, 1000000)

	name := NameFlag
	genDoc.ChainID = name

	// overwrite account
	acc := genDoc.Accounts[0]
	acc.Address = genDoc.Validators[0].UnbondTo[0].Address
	genDoc.Accounts[0] = acc

	buf, buf2, n, err := new(bytes.Buffer), new(bytes.Buffer), new(int64), new(error)
	binary.WriteJSON(genDoc, buf, n, err)
	IfExit(*err)
	IfExit(json.Indent(buf2, buf.Bytes(), "", "\t"))
	genesisBytes := buf2.Bytes()

	// create directory to save priv validators and genesis.json
	IfExit(os.MkdirAll(dir, 0700))
	IfExit(ioutil.WriteFile(path.Join(dir, "genesis.json"), genesisBytes, 0644))
	v := validators[0]
	buf, n, err = new(bytes.Buffer), new(int64), new(error)
	binary.WriteJSON(v, buf, n, err)
	IfExit(*err)
	valBytes := buf.Bytes()
	IfExit(ioutil.WriteFile(path.Join(dir, "priv_validator.json"), valBytes, 0600))
	IfExit(ioutil.WriteFile(path.Join(dir, "config.toml"), []byte(defaultConfig), 0600))
	fmt.Printf("genesis.json, config.toml and priv_validator.json files saved in %s\n", dir)
}

func cliGenesis(cmd *cobra.Command, args []string) {
	if cmd.Flags().Lookup("single").Changed {
		if len(args) != 1 {
			Exit(fmt.Errorf("Enter a directory to save the genesis.json and priv_validator.json to"))
		}

		singleUserChain(args[0])
		return
	}

	if len(args) != 1 {
		Exit(fmt.Errorf("Enter the number of validators to create"))
	}
	N, err := strconv.Atoi(args[0])
	if err != nil {
		Exit(fmt.Errorf("Please provide an integer number of validators to create"))
	}

	fmt.Println("Generating accounts ...")
	genDoc, _, validators := state.RandGenesisDoc(N, true, 100000, N, false, 1000)

	name := NameFlag
	genDoc.ChainID = name

	buf, buf2, n := new(bytes.Buffer), new(bytes.Buffer), new(int64)
	binary.WriteJSON(genDoc, buf, n, &err)
	IfExit(err)
	IfExit(json.Indent(buf2, buf.Bytes(), "", "\t"))
	genesisBytes := buf2.Bytes()

	// create directory to save priv validators and genesis.json
	d := path.Join(DataContainersPath, name)
	IfExit(os.MkdirAll(d, 0700))
	IfExit(ioutil.WriteFile(path.Join(d, "genesis.json"), genesisBytes, 0644))
	for i, v := range validators {
		buf, n = new(bytes.Buffer), new(int64)
		binary.WriteJSON(v, buf, n, &err)
		IfExit(err)
		valBytes := buf.Bytes()
		IfExit(ioutil.WriteFile(path.Join(d, fmt.Sprintf("priv_validator_%d.json", i)), valBytes, 0600))
	}
	fmt.Printf("genesis.json and priv_validator_X.json files saved in %s\n", d)
}
