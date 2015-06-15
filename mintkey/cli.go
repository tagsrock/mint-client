package main

import (
	"encoding/hex"
	"fmt"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/eris-keys/crypto"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/binary"
)

func cliConvertAddressToPrivValidator(cmd *cobra.Command, args []string) {
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

	privVal := struct {
		Address    []byte                 `json:"address"`
		PubKey     account.PubKeyEd25519  `json:"pub_key"`
		PrivKey    account.PrivKeyEd25519 `json:"priv_key"`
		LastHeight int                    `json:"last_height"`
		LastRound  int                    `json:"last_round"`
		LastStep   int                    `json:"last_step"`
	}{
		Address: addrBytes,
		PubKey:  account.PubKeyEd25519(pub),
		PrivKey: account.PrivKeyEd25519(key.PrivateKey),
	}

	fmt.Println(string(binary.JSONBytes(privVal)))
}
