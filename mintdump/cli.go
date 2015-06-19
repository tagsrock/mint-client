package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	acm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/binary"
	cfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
	dbm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/db"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/merkle"
	sm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
)

//------------------------------------------------------------------------------
// core dump/restore functions

// dump the latest state to json
func CoreDump() []byte {
	// Get State
	stateDB := dbm.GetDB("state")
	st := sm.LoadState(stateDB)

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

		storage := merkle.NewIAVLTree(binary.BasicCodec, binary.BasicCodec, 1024, stateDB)
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
	binary.WriteJSON(stJ, w, n, err)
	ifExit(*err)
	w2 := new(bytes.Buffer)
	json.Indent(w2, w.Bytes(), "", "\t")
	return w2.Bytes()
}

// restore state from json blob
// set tendermint config before calling
func CoreRestore(chainID string, jsonBytes []byte) {
	var stJ State
	var err error
	binary.ReadJSON(&stJ, jsonBytes, &err)
	ifExit(err)

	st := new(sm.State)

	st.ChainID = chainID
	st.BondedValidators = stJ.BondedValidators
	st.LastBondedValidators = stJ.LastBondedValidators
	st.UnbondingValidators = stJ.UnbondingValidators

	stateDB := dbm.GetDB("state")

	// fill the accounts tree
	accounts := merkle.NewIAVLTree(binary.BasicCodec, acm.AccountCodec, 1000, stateDB)
	for _, account := range stJ.Accounts {
		accounts.Set(account.Address, account.Copy())
	}

	// fill the storage tree for each contract
	for _, accStorage := range stJ.AccountsStorage {
		st := merkle.NewIAVLTree(binary.BasicCodec, binary.BasicCodec, 1024, stateDB)
		for _, accSt := range accStorage.Storage {
			set := st.Set(accSt.Key, accSt.Value)
			if !set {
				panic("failed to update storage tree")
			}
		}
		// TODO: sanity check vs acc.StorageRoot

		st.Save()
	}

	valInfos := merkle.NewIAVLTree(binary.BasicCodec, sm.ValidatorInfoCodec, 0, stateDB)
	for _, valInfo := range stJ.ValidatorInfos {
		valInfos.Set(valInfo.Address, valInfo)
	}

	nameReg := merkle.NewIAVLTree(binary.BasicCodec, sm.NameRegCodec, 0, stateDB)
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
	if len(args) != 0 {
		exit(fmt.Errorf("Enter the path to a json file containing a state dump, followed by the new chain id"))
	}
	file, chainID := args[0], args[1]

	b, err := ioutil.ReadFile(file)
	ifExit(err)

	CoreRestore(chainID, b)

	stateDB := dbm.GetDB("state")
	newState := sm.LoadState(stateDB)
	fmt.Printf("State hash: %X\n", newState.Hash())

}

func cliDump(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		config.Set("db_dir", args[0])
		cfg.ApplyConfig(config) // Notify modules of new config
	}

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
