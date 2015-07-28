package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	acm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	dbm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/db"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/merkle"
	sm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

//------------------------------------------------------------------------------
// core dump/restore functions

// dump the latest state to json
func CoreDump() []byte {
	// Get State
	stateDB := dbm.GetDB("state")
	st := sm.LoadState(stateDB)
	if st == nil {
		Exit(fmt.Errorf("Error: state loaded from %s is nil!", config.GetString("db_dir")))
	}

	stJ := new(State)
	stJ.BondedValidators = st.BondedValidators
	stJ.LastBondedValidators = st.LastBondedValidators
	stJ.UnbondingValidators = st.UnbondingValidators

	// iterate through accounts tree
	// track storage roots as we go
	storageRoots := [][]byte{}
	st.GetAccounts().Iterate(func(key interface{}, value interface{}) (stopped bool) {
		acc := value.(*acm.Account)
		stJ.Accounts = append(stJ.Accounts, acc)
		storageRoots = append(storageRoots, acc.StorageRoot)
		return false
	})

	// grab all storage
	for i, root := range storageRoots {
		if len(root) == 0 {
			continue
		}
		accStorage := &AccountStorage{Address: stJ.Accounts[i].Address}

		storage := merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, 1024, stateDB)
		storage.Load(root)
		storage.Iterate(func(key interface{}, value interface{}) (stopped bool) {
			k, v := key.([]byte), value.([]byte)
			accStorage.Storage = append(accStorage.Storage, &Storage{k, v})
			return false
		})
		stJ.AccountsStorage = append(stJ.AccountsStorage, accStorage)
	}

	// get all validator infos
	st.GetValidatorInfos().Iterate(func(key interface{}, value interface{}) (stopped bool) {
		vi := value.(*sm.ValidatorInfo)
		stJ.ValidatorInfos = append(stJ.ValidatorInfos, vi)
		return false
	})

	// get all name entries
	st.GetNames().Iterate(func(key interface{}, value interface{}) (stopped bool) {
		name := value.(*types.NameRegEntry)
		stJ.NameReg = append(stJ.NameReg, name)
		return false
	})

	w, n, err := new(bytes.Buffer), new(int64), new(error)
	wire.WriteJSON(stJ, w, n, err)
	IfExit(*err)
	w2 := new(bytes.Buffer)
	json.Indent(w2, w.Bytes(), "", "\t")
	return w2.Bytes()
}

// restore state from json blob
// set tendermint config before calling
func CoreRestore(chainID string, jsonBytes []byte) {
	var stJ State
	var err error
	wire.ReadJSON(&stJ, jsonBytes, &err)
	IfExit(err)

	st := new(sm.State)

	st.ChainID = chainID
	st.BondedValidators = stJ.BondedValidators
	st.LastBondedValidators = stJ.LastBondedValidators
	st.UnbondingValidators = stJ.UnbondingValidators

	stateDB := dbm.GetDB("state")

	// fill the accounts tree
	accounts := merkle.NewIAVLTree(wire.BasicCodec, acm.AccountCodec, 1000, stateDB)
	for _, account := range stJ.Accounts {
		accounts.Set(account.Address, account.Copy())
	}

	// fill the storage tree for each contract
	for _, accStorage := range stJ.AccountsStorage {
		st := merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, 1024, stateDB)
		for _, accSt := range accStorage.Storage {
			set := st.Set(accSt.Key, accSt.Value)
			if !set {
				panic("failed to update storage tree")
			}
		}
		// TODO: sanity check vs acc.StorageRoot

		st.Save()
	}

	valInfos := merkle.NewIAVLTree(wire.BasicCodec, sm.ValidatorInfoCodec, 0, stateDB)
	for _, valInfo := range stJ.ValidatorInfos {
		valInfos.Set(valInfo.Address, valInfo)
	}

	nameReg := merkle.NewIAVLTree(wire.BasicCodec, sm.NameRegCodec, 0, stateDB)
	for _, entry := range stJ.NameReg {
		nameReg.Set(entry.Name, entry)
	}

	// persists accounts/valInfos/nameReg trees
	st.SetAccounts(accounts)
	st.SetValidatorInfos(valInfos)
	st.SetNameReg(nameReg)
	st.SetDB(stateDB)
	st.Save()
}

//------------------------------------------------------------------------------
// cli wrappers

func cliRestore(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		Exit(fmt.Errorf("Enter the chain id"))
	}
	chainID := args[0]

	fi, _ := os.Stdin.Stat()
	if fi.Size() == 0 {
		Exit(fmt.Errorf("Please pass data to restore on Stdin"))
	}
	b, err := ioutil.ReadAll(os.Stdin)
	IfExit(err)

	CoreRestore(chainID, b)

	stateDB := dbm.GetDB("state")
	newState := sm.LoadState(stateDB)
	fmt.Printf("State hash: %X\n", newState.Hash())

}

func cliDump(cmd *cobra.Command, args []string) {
	fmt.Println(string(CoreDump()))
}

//------------------------------------------------------------------------------
// types

type State struct {
	BondedValidators     *sm.ValidatorSet      `json:"bonded_validators"`
	LastBondedValidators *sm.ValidatorSet      `json:"last_bonded_validators"`
	UnbondingValidators  *sm.ValidatorSet      `json:"unbonding_validators"`
	Accounts             []*acm.Account        `json:"accounts"`
	AccountsStorage      []*AccountStorage     `json:"accounts_storage"`
	ValidatorInfos       []*sm.ValidatorInfo   `json:"validator_infos"`
	NameReg              []*types.NameRegEntry `json:"namereg"`
}

type AccountStorage struct {
	Address []byte     `json:"address"`
	Storage []*Storage `json:"storage"`
}

type Storage struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}
