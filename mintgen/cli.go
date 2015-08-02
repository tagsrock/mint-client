package main

import (
	"bytes"
	"encoding/csv"
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

	for i, v := range validators {
		buf, n = new(bytes.Buffer), new(int64)
		wire.WriteJSON(v, buf, n, &err)
		IfExit(err)
		valBytes := buf.Bytes()
		if len(validators) > 1 {
			mulDir := fmt.Sprintf("%s_%d", DirFlag, i)
			IfExit(os.MkdirAll(mulDir, 0700))
			IfExit(ioutil.WriteFile(path.Join(mulDir, "priv_validator.json"), valBytes, 0600))
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
	//TODO better logic -> refactor; set defaults, or add them as flags??
	if PubkeyFlag == "" && CsvPathFlag == "" {
		Exit(fmt.Errorf("Enter one or more pub keys or path to .csv"))
	}

	var pubkeys []string
	var amts []string
	var names []string
	//	var perms []string

	if CsvPathFlag != "" {
		//TODO figure out perms
		pubkeys, amts, names, _ = parseCsv(CsvPathFlag)

	} else {
		pubkeys = strings.Split(PubkeyFlag, " ")
	}

	pubKeyBytes := make([][]byte, len(pubkeys))
	for i, k := range pubkeys {
		var err error
		pubKeyBytes[i], err = hex.DecodeString(k)
		if err != nil {
			Exit(fmt.Errorf("Pubkey (%s) is invalid hex: %v", k, err))
		}
	}

	amt := make([]int64, len(amts))
	for i, a := range amts {
		var err error
		amt[i], err = strconv.ParseInt(a, 10, 64)
		if err != nil {
			Exit(fmt.Errorf("Invalid amount: %v", err))
		}
	}

	genDoc := state.GenesisDoc{
		ChainID: chainID,
	}

	genDoc.Accounts = make([]state.GenesisAccount, len(pubkeys))
	genDoc.Validators = make([]state.GenesisValidator, len(pubkeys))

	var pubKey account.PubKeyEd25519
	unbAmt := int64(1) << 50

	for i, kb := range pubKeyBytes {
		copy(pubKey[:], kb)
		addr := pubKey.Address()

		genDoc.Accounts[i] = state.GenesisAccount{
			Address: addr,
			Amount:  amt[i],
			Name:    names[i],
		}
		genDoc.Validators[i] = state.GenesisValidator{
			PubKey: pubKey,
			Amount: amt[i],
			UnbondTo: []state.BasicAccount{
				state.BasicAccount{
					Address: addr,
					Amount:  unbAmt,
				},
			},
		}
	}

	buf, buf2, n, err := new(bytes.Buffer), new(bytes.Buffer), new(int64), new(error)
	wire.WriteJSON(genDoc, buf, n, err)
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

//takes a csv in the format defined [here]
func parseCsv(path string) (pubkeys, amts, names, perms []string) {

	csvFile, err := os.Open(path)
	if err != nil {
		Exit(fmt.Errorf("Couldn't open file: %s: %v", path, err))
	}

	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	//r.FieldsPerRecord = # of records expected
	params, err := r.ReadAll()
	if err != nil {
		Exit(fmt.Errorf("Couldn't read file: %v", err))

	}

	pubkeys = make([]string, len(params))
	amts = make([]string, len(params))
	names = make([]string, len(params))
	perms = make([]string, len(params))
	for i, each := range params {
		pubkeys[i] = each[0]
		amts[i] = each[1]
		names[i] = each[2]
		perms[i] = each[3]

	}
	return pubkeys, amts, names, perms
}
