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
// mintgen cli

func cliKnown(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		Exit(fmt.Errorf("Enter a chain_id"))
	}
	chainID := args[0]

	var genDoc *state.GenesisDoc
	var err error
	if CsvPathFlag != "" {
		//TODO figure out perms
		pubkeys, amts, names, _ := parseCsv(CsvPathFlag)

		// convert amts to ints
		amt := make([]int64, len(amts))
		for i, a := range amts {
			amt[i], err = strconv.ParseInt(a, 10, 64)
			if err != nil {
				Exit(fmt.Errorf("Invalid amount: %v", err))
			}
		}

		// convert pubkey hex strings to struct
		pubKeys := pubKeyStringsToPubKeys(pubkeys)

		genDoc = newGenDoc(chainID, len(pubKeys), len(pubKeys))
		for i, pk := range pubKeys {
			genDocAddAccountAndValidator(genDoc, pk, amt[i], names[i], i)
		}

	} else if PubkeyFlag != "" {
		pubkeys := strings.Split(PubkeyFlag, " ")
		amt := int64(1) << 50
		pubKeys := pubKeyStringsToPubKeys(pubkeys)

		genDoc = newGenDoc(chainID, len(pubKeys), len(pubKeys))
		for i, pk := range pubKeys {
			genDocAddAccountAndValidator(genDoc, pk, amt, "", i)
		}

	} else {
		privJSON := readStdinTimeout()
		genDoc = genesisFromPrivValBytes(chainID, privJSON)
	}

	buf, buf2, n := new(bytes.Buffer), new(bytes.Buffer), new(int64)
	wire.WriteJSON(genDoc, buf, n, &err)
	IfExit(err)
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

	for i, v := range validators {
		buf, n = new(bytes.Buffer), new(int64)
		wire.WriteJSON(v, buf, n, &err)
		IfExit(err)
		valBytes := buf.Bytes()
		if len(validators) > 1 {
			mulDir := fmt.Sprintf("%s_%d", DirFlag, i)
			IfExit(os.MkdirAll(mulDir, 0700))
			IfExit(ioutil.WriteFile(path.Join(mulDir, "priv_validator.json"), valBytes, 0600))
			IfExit(ioutil.WriteFile(path.Join(mulDir, "config.toml"), []byte(setDefaultConfig(i, chainID, SeedsFlag)), 0644))
			IfExit(ioutil.WriteFile(path.Join(mulDir, "genesis.json"), genesisBytes, 0644))
		} else {
			IfExit(ioutil.WriteFile(path.Join(DirFlag, "priv_validator.json"), valBytes, 0600))
			IfExit(ioutil.WriteFile(path.Join(DirFlag, "config.toml"), []byte(setDefaultConfig(i, chainID, SeedsFlag)), 0644))
			IfExit(ioutil.WriteFile(path.Join(DirFlag, "genesis.json"), genesisBytes, 0644))
		}
	}
	fmt.Printf("config.toml, genesis.json and priv_validator.json files saved in %s\n", DirFlag)
}

//-----------------------------------------------------------------------------
// gendoc convenience functions

func newGenDoc(chainID string, nVal, nAcc int) *state.GenesisDoc {
	genDoc := state.GenesisDoc{
		ChainID: chainID,
	}
	genDoc.Accounts = make([]state.GenesisAccount, nAcc)
	genDoc.Validators = make([]state.GenesisValidator, nVal)
	return &genDoc
}

// genesis file with only one validator, using priv_validator.json
func genesisFromPrivValBytes(chainID string, privJSON []byte) *state.GenesisDoc {
	var err error
	privVal := wire.ReadJSON(&state.PrivValidator{}, privJSON, &err).(*state.PrivValidator)
	if err != nil {
		Exit(fmt.Errorf("Error reading PrivValidator on stdin: %v\n", err))
	}
	pubKey := privVal.PubKey
	amt := int64(1) << 50

	genDoc := newGenDoc(chainID, 1, 1)

	genDocAddAccountAndValidator(genDoc, pubKey, amt, "", 0)

	return genDoc
}

func genDocAddAccountAndValidator(genDoc *state.GenesisDoc, pubKey account.PubKeyEd25519, amt int64, name string, index int) {
	addr := pubKey.Address()
	genDoc.Accounts[index] = state.GenesisAccount{
		Address: addr,
		Amount:  amt,
		Name:    name,
	}
	genDoc.Validators[index] = state.GenesisValidator{
		PubKey: pubKey,
		Amount: amt,
		Name:   name,
		UnbondTo: []state.BasicAccount{
			state.BasicAccount{
				Address: addr,
				Amount:  amt,
			},
		},
	}
}

//-----------------------------------------------------------------------------
// util functions

func setDefaultConfig(num int, mon, seeds string) []byte {
	//build moniker
	moniker := fmt.Sprintf("%s_%d", mon, num)
	var defaultConfig = fmt.Sprintf(`
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

moniker = "%s"
node_laddr = "0.0.0.0:46656"
seeds = "%s"
fast_sync = false
db_backend = "leveldb"
log_level = "debug"
rpc_laddr = "0.0.0.0:46657"
`, moniker, seeds)

	return []byte(defaultConfig)
}

// convert hex strings to ed25519 pubkeys
func pubKeyStringsToPubKeys(pubkeys []string) []account.PubKeyEd25519 {
	pubKeys := make([]account.PubKeyEd25519, len(pubkeys))
	for i, k := range pubkeys {
		pubBytes, err := hex.DecodeString(k)
		if err != nil {
			Exit(fmt.Errorf("Pubkey (%s) is invalid hex: %v", k, err))
		}
		copy(pubKeys[i][:], pubBytes)
	}
	return pubKeys
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

const stdinTimeoutSeconds = 1

// read the priv validator json off stdin or timeout and fail
func readStdinTimeout() []byte {
	ch := make(chan []byte, 1)
	go func() {
		privJSON, err := ioutil.ReadAll(os.Stdin)
		IfExit(err)
		ch <- privJSON
	}()
	ticker := time.Tick(time.Second * stdinTimeoutSeconds)
	select {
	case <-ticker:
		Exit(fmt.Errorf("Please pass a priv_validator.json on stdin, or specify either a pubkey with --pub or csv file with --csv"))
	case privJSON := <-ch:
		return privJSON
	}
	return nil
}
