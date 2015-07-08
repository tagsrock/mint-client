package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
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

const stdinTimeoutSeconds = 1

func cliSingle(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		Exit(fmt.Errorf("Enter a chain_id"))
	}

	chainID := args[0]

	var pubKey account.PubKeyEd25519
	var addr []byte
	if PubkeyFlag == "" {
		// block reading on stdin, wait 1 second
		ch := make(chan []byte, 1)
		go func() {
			privJSON, err := ioutil.ReadAll(os.Stdin)
			IfExit(err)
			ch <- privJSON
		}()
		ticker := time.Tick(time.Second * stdinTimeoutSeconds)
		select {
		case <-ticker:
			Exit(fmt.Errorf("Please pass a priv_validator.json on stdin or specify a pubkey with --pub"))
		case privJSON := <-ch:
			var err error
			privVal := binary.ReadJSON(&state.PrivValidator{}, privJSON, &err).(*state.PrivValidator)
			if err != nil {
				Exit(fmt.Errorf("Error reading PrivValidator on stdin: %v\n", err))
			}
			pubKey = privVal.PubKey
			addr = privVal.Address
		}
	} else {
		pubKeyBytes, err := hex.DecodeString(PubkeyFlag)
		if err != nil {
			Exit(fmt.Errorf("Pubkey (%s) is invalid hex: %v", PubkeyFlag, err))
		}
		pubKey = account.PubKeyEd25519(pubKeyBytes)
		addr = pubKey.Address()
	}

	amt := uint64(1) << 61
	//  build gendoc
	genDoc := state.GenesisDoc{
		ChainID: chainID,
		Accounts: []state.GenesisAccount{
			state.GenesisAccount{
				Address: addr,
				Amount:  amt,
			},
		},
		Validators: []state.GenesisValidator{
			state.GenesisValidator{
				PubKey: pubKey,
				Amount: amt,
				UnbondTo: []state.BasicAccount{
					state.BasicAccount{
						Address: addr,
						Amount:  amt,
					},
				},
			},
		},
	}

	buf, buf2, n, err := new(bytes.Buffer), new(bytes.Buffer), new(int64), new(error)
	binary.WriteJSON(genDoc, buf, n, err)
	IfExit(*err)
	IfExit(json.Indent(buf2, buf.Bytes(), "", "\t"))
	genesisBytes := buf2.Bytes()

	fmt.Println(string(genesisBytes))
}

func cliRandom(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		Exit(fmt.Errorf("Enter the number of validators and a chain_id"))
	}
	N, err := strconv.Atoi(args[0])
	if err != nil {
		Exit(fmt.Errorf("Please provide an integer number of validators to create"))
	}

	chainID := args[1]

	fmt.Println("Generating accounts ...")
	genDoc, _, validators := state.RandGenesisDoc(N, true, 100000, N, false, 1000)

	genDoc.ChainID = chainID

	buf, buf2, n := new(bytes.Buffer), new(bytes.Buffer), new(int64)
	binary.WriteJSON(genDoc, buf, n, &err)
	IfExit(err)
	IfExit(json.Indent(buf2, buf.Bytes(), "", "\t"))
	genesisBytes := buf2.Bytes()

	// create directory to save priv validators and genesis.json
	if DirFlag == "" {
		DirFlag = path.Join(DataContainersPath, chainID)
	}
	if _, err := os.Stat(DirFlag); err != nil {
		IfExit(os.MkdirAll(DirFlag, 0700))
	}
	// XXX: we're gonna write to it even if it already exists
	IfExit(ioutil.WriteFile(path.Join(DirFlag, "genesis.json"), genesisBytes, 0644))
	for i, v := range validators {
		buf, n = new(bytes.Buffer), new(int64)
		binary.WriteJSON(v, buf, n, &err)
		IfExit(err)
		valBytes := buf.Bytes()
		IfExit(ioutil.WriteFile(path.Join(DirFlag, fmt.Sprintf("priv_validator_%d.json", i)), valBytes, 0600))
	}
	fmt.Printf("genesis.json and priv_validator_X.json files saved in %s\n", DirFlag)
}
