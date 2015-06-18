package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/eris-keys/crypto"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/binary"
)

type PrivValidator struct {
	Address    []byte                 `json:"address"`
	PubKey     account.PubKeyEd25519  `json:"pub_key"`
	PrivKey    account.PrivKeyEd25519 `json:"priv_key"`
	LastHeight int                    `json:"last_height"`
	LastRound  int                    `json:"last_round"`
	LastStep   int                    `json:"last_step"`
}

func cliConvertErisKeyToPrivValidator(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) == 0 {
		exit(fmt.Errorf("Please enter the address of your key"))
	}

	addr := args[0]
	addrBytes, err := hex.DecodeString(addr)
	ifExit(err)

	keyStore := crypto.NewKeyStorePlain(DefaultKeyStore)
	key, err := keyStore.GetKey(addrBytes, "")
	ifExit(err)

	pub, err := key.Pubkey()
	ifExit(err)

	privVal := PrivValidator{
		Address: addrBytes,
		PubKey:  account.PubKeyEd25519(pub),
		PrivKey: account.PrivKeyEd25519(key.PrivateKey),
	}

	fmt.Println(string(binary.JSONBytes(privVal)))
}

func cliConvertPrivValidatorToErisKey(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) == 0 {
		exit(fmt.Errorf("Please enter the path to the priv_validator.json"))
	}

	pvf := args[0]
	b, err := ioutil.ReadFile(pvf)
	ifExit(err)

	pv := new(PrivValidator)
	binary.ReadJSON(pv, b, &err)
	ifExit(err)

	keyStore := crypto.NewKeyStorePlain(DefaultKeyStore)

	key := &crypto.Key{
		Id:         uuid.NewRandom(),
		Type:       crypto.KeyTypeEd25519,
		Address:    pv.Address,
		PrivateKey: pv.PrivKey,
	}

	fmt.Printf("Converted key for address %X\n", key.Address)
	fmt.Printf("Storing %X in keyStore (%s)\n", key.Address, DefaultKeyStore)

	ifExit(keyStore.StoreKey(key, ""))
}
