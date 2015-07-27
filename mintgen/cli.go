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
	"strings"
	"time"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
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
			privVal := wire.ReadJSON(&state.PrivValidator{}, privJSON, &err).(*state.PrivValidator)
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
		copy(pubKey[:], pubKeyBytes)
		addr = pubKey.Address()
	}

	amt := int64(1) << 60
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
	wire.WriteJSON(genDoc, buf, n, err)
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

	// RandGenesisDoc produces random accounts and validators.
	// Give the validators accounts:
	genDoc.Accounts = make([]state.GenesisAccount, N)
	for i, pv := range validators {
		genDoc.Accounts[i] = state.GenesisAccount{
			Address: pv.Address,
			Amount:  int64(2) << 50,
		}
	}

	buf, buf2, n := new(bytes.Buffer), new(bytes.Buffer), new(int64)
	wire.WriteJSON(genDoc, buf, n, &err)
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
	IfExit(ioutil.WriteFile(path.Join(DirFlag, "config.toml"), []byte(defaultConfig), 0644))
	IfExit(ioutil.WriteFile(path.Join(DirFlag, "genesis.json"), genesisBytes, 0644))
	fmt.Println("dirflag", DirFlag)
	for i, v := range validators {
		buf, n = new(bytes.Buffer), new(int64)
		wire.WriteJSON(v, buf, n, &err)
		IfExit(err)
		valBytes := buf.Bytes()
		if len(validators) > 1 {
			mulDir := DirFlag + "_" + strconv.Itoa(i)
			fmt.Println("muldir", mulDir)
			fmt.Println("i", i)
			IfExit(os.MkdirAll(mulDir, 0700))
			IfExit(ioutil.WriteFile(path.Join(mulDir, fmt.Sprintf("priv_validator_%d.json", i)), valBytes, 0600))
			IfExit(ioutil.WriteFile(path.Join(mulDir, "config.toml"), []byte(defaultConfig), 0644))
			IfExit(ioutil.WriteFile(path.Join(mulDir, "genesis.json"), genesisBytes, 0644))
		} else {
			IfExit(ioutil.WriteFile(path.Join(DirFlag, "priv_validator.json"), valBytes, 0600))
		}
	}
	fmt.Printf("config.toml, genesis.json and priv_validator_X.json files saved in %s\n", DirFlag)
}

func cliMulti(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		Exit(fmt.Errorf("Enter a chain_id"))
	}
	chainID := args[0]

	//TODO convert to addrs
	if PubkeyFlag == "" {
		Exit(fmt.Errorf("Enter one or more pub keys"))
	}

	pubkeys := strings.Split(PubkeyFlag, " ")

	pubKeyBytes := make([][]byte, len(pubkeys))
	for i, k := range pubkeys {
		var err error
		pubKeyBytes[i], err = hex.DecodeString(k)
		if err != nil {
			Exit(fmt.Errorf("Pubkey (%s) is invalid hex: %v", k, err))
		}
	}

	addrs := make([][]byte, len(pubkeys))
	amt := int64(1) << 50

	genDoc := state.GenesisDoc{
		ChainID: chainID,
	}

	genDoc.Accounts = make([]state.GenesisAccount, len(addrs))
	genDoc.Validators = make([]state.GenesisValidator, len(addrs))

	for i, kb := range pubKeyBytes {
		pubKey := account.PubKeyEd25519(kb)
		addr := pubKey.Address()

		genDoc.Accounts[i] = state.GenesisAccount{
			Address: addr,
			Amount:  amt,
		}
		genDoc.Validators[i] = state.GenesisValidator{
			PubKey: pubKey,
			Amount: amt,
			UnbondTo: []state.BasicAccount{
				state.BasicAccount{
					Address: addr,
					Amount:  amt,
				},
			},
		}
	}

	buf, buf2, n, err := new(bytes.Buffer), new(bytes.Buffer), new(int64), new(error)
	binary.WriteJSON(genDoc, buf, n, err)
	IfExit(*err)
	IfExit(json.Indent(buf2, buf.Bytes(), "", "\t"))
	genesisBytes := buf2.Bytes()

	fmt.Println(string(genesisBytes))
	if DirFlag == "" {
		DirFlag = path.Join(DataContainersPath, chainID)
	}
	if _, err := os.Stat(DirFlag); err != nil {
		IfExit(os.MkdirAll(DirFlag, 0700))
	}

	IfExit(ioutil.WriteFile(path.Join(DirFlag, "genesis.json"), genesisBytes, 0644))
	fmt.Printf("genesis.json saved in %s\n", DirFlag)
}
