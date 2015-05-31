package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/tendermint/tendermint/account"
)

// do we really need these?
func cliInput(c *cli.Context) {
	pubkey, amtS, nonceS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("addr")
	input, err := coreInput(pubkey, amtS, nonceS, addr)
	ifExit(err)
	fmt.Printf("%s\n", input)
}

func cliOutput(c *cli.Context) {
	addr, amtS := c.String("addr"), c.String("amt")
	output, err := coreOutput(addr, amtS)
	ifExit(err)
	fmt.Printf("%s\n", output)
}

func cliSend(c *cli.Context) {
	chainID := c.String("chainID")
	pubkey, amtS, nonceS, addr, toAddr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("addr"), c.String("to")
	tx, err := coreSend(chainID, pubkey, addr, toAddr, amtS, nonceS)
	ifExit(err)
	fmt.Printf("%v\n", tx)

	sign, broadcast := c.String("sign"), c.String("broadcast")
	if sign != "" {
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Inputs[0].Address)
		sig, err := coreSign(signBytes, addrHex, sign)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Inputs[0].Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast != "" {
		receipt, err := coreBroadcast(tx, broadcast)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliName(c *cli.Context) {
	chainID := c.String("chainID")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")

	name, data := c.String("name"), c.String("data")
	tx, err := coreName(chainID, pubkey, addr, amtS, nonceS, feeS, name, data)
	ifExit(err)
	fmt.Printf("%v\n", tx)
	sign, broadcast := c.String("sign"), c.String("broadcast")
	if sign != "" {
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Input.Address)
		sig, err := coreSign(signBytes, addrHex, sign)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Input.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast != "" {
		receipt, err := coreBroadcast(tx, broadcast)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliCall(c *cli.Context) {
	chainID := c.String("chainID")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")

	toAddr, gasS, data := c.String("name"), c.String("gas"), c.String("data")
	tx, err := coreCall(chainID, pubkey, addr, toAddr, amtS, nonceS, gasS, feeS, data)
	ifExit(err)
	fmt.Printf("%v\n", tx)
	sign, broadcast := c.String("sign"), c.String("broadcast")
	if sign != "" {
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Input.Address)
		sig, err := coreSign(signBytes, addrHex, sign)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Input.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast != "" {
		receipt, err := coreBroadcast(tx, broadcast)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}
