package main

import (
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
)

// do we really need these?
/*
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
*/
func cliSend(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, amtS, nonceS, addr, toAddr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("addr"), c.String("to")
	tx, err := coreSend(nodeAddr, pubkey, addr, toAddr, amtS, nonceS)
	ifExit(err)
	fmt.Printf("%v\n", tx)

	sign, broadcast := c.Bool("sign"), c.Bool("broadcast")
	if sign {
		signAddr := c.String("sign-addr")
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Inputs[0].Address)
		sig, err := coreSign(signBytes, addrHex, signAddr)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Inputs[0].Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast {
		receipt, err := coreBroadcast(tx, nodeAddr)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliName(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")

	if c.IsSet("data") && c.IsSet("data-file") {
		exit(fmt.Errorf("Please specify only one of --data and --data-file"))
	}
	name, data, dataFile := c.String("name"), c.String("data"), c.String("data-file")
	if data == "" && dataFile != "" {
		b, err := ioutil.ReadFile(dataFile)
		ifExit(err)
		data = string(b)
	}
	tx, err := coreName(nodeAddr, pubkey, addr, amtS, nonceS, feeS, name, data)
	ifExit(err)
	fmt.Printf("%v\n", tx)
	sign, broadcast := c.Bool("sign"), c.Bool("broadcast")
	if sign {
		signAddr := c.String("sign-addr")
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Input.Address)
		sig, err := coreSign(signBytes, addrHex, signAddr)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Input.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast {
		receipt, err := coreBroadcast(tx, nodeAddr)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliCall(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")

	toAddr, gasS, data := c.String("to"), c.String("gas"), c.String("data")
	tx, err := coreCall(nodeAddr, pubkey, addr, toAddr, amtS, nonceS, gasS, feeS, data)
	ifExit(err)
	fmt.Printf("%v\n", tx)
	sign, broadcast := c.Bool("sign"), c.Bool("broadcast")
	if sign {
		signAddr := c.String("sign-addr")
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Input.Address)
		sig, err := coreSign(signBytes, addrHex, signAddr)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Input.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast {
		receipt, err := coreBroadcast(tx, nodeAddr)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliBond(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, amtS, nonceS, unbondAddr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("unbond-to")
	tx, err := coreBond(nodeAddr, pubkey, unbondAddr, amtS, nonceS)
	ifExit(err)
	fmt.Printf("%v\n", tx)

	sign, broadcast := c.Bool("sign"), c.Bool("broadcast")
	if sign {
		signAddr := c.String("sign-addr")
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Inputs[0].Address)
		sig, err := coreSign(signBytes, addrHex, signAddr)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		// for now we expect the bonder to send from
		// the same account as his key, but this will obviously need to change
		tx.Inputs[0].Signature = sigED
		tx.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast {
		receipt, err := coreBroadcast(tx, nodeAddr)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliUnbond(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	addr, height := c.String("addr"), c.String("height")

	tx, err := coreUnbond(addr, height)
	ifExit(err)
	fmt.Printf("%v\n", tx)

	sign, broadcast := c.Bool("sign"), c.Bool("broadcast")
	if sign {
		signAddr := c.String("sign-addr")
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Address)
		sig, err := coreSign(signBytes, addrHex, signAddr)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast {
		receipt, err := coreBroadcast(tx, nodeAddr)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}

func cliRebond(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	addr, height := c.String("addr"), c.String("height")

	tx, err := coreRebond(addr, height)
	ifExit(err)
	fmt.Printf("%v\n", tx)

	sign, broadcast := c.Bool("sign"), c.Bool("broadcast")
	if sign {
		signAddr := c.String("sign-addr")
		signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx))
		addrHex := fmt.Sprintf("%X", tx.Address)
		sig, err := coreSign(signBytes, addrHex, signAddr)
		ifExit(err)
		sigED := account.SignatureEd25519(sig[:])
		tx.Signature = sigED
		fmt.Printf("%X\n", sig)
	}
	if broadcast {
		receipt, err := coreBroadcast(tx, nodeAddr)
		ifExit(err)
		fmt.Printf("%X\n", receipt.TxHash)
	}
}
