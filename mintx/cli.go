package main

import (
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/codegangsta/cli"
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
	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
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
	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
}

func cliCall(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")
	toAddr, gasS, data := c.String("to"), c.String("gas"), c.String("data")
	tx, err := coreCall(nodeAddr, pubkey, addr, toAddr, amtS, nonceS, gasS, feeS, data)
	ifExit(err)
	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
}

func cliPermissions(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, nonceS, addr := c.String("pubkey"), c.String("nonce"), c.String("addr")

	// all functions take at least 2 args (+ name)
	if len(c.Args()) < 3 {
		ifExit(fmt.Errorf("Please enter the permission function you'd like to call, followed by it's arguments"))
	}
	permFunc := c.Args()[0]
	tx, err := corePermissions(nodeAddr, pubkey, addr, nonceS, permFunc, c.Args()[1:])
	ifExit(err)
	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
}

func cliNewAccount(c *cli.Context) {
	/*
		chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
		pubkey := c.String("pubkey")

		tx, err := coreNewAccount(nodeAddr, pubkey, chainID)
		ifExit(err)

		logger.Debugf("%v\n", tx)
		signAndBroadcast(c, chainID, nodeAddr, tx)
	*/
}

func cliBond(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	pubkey, amtS, nonceS, unbondAddr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("unbond-to")
	tx, err := coreBond(nodeAddr, pubkey, unbondAddr, amtS, nonceS)
	ifExit(err)

	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
}

func cliUnbond(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	addr, height := c.String("addr"), c.String("height")
	tx, err := coreUnbond(addr, height)
	ifExit(err)
	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
}

func cliRebond(c *cli.Context) {
	chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
	addr, height := c.String("addr"), c.String("height")
	tx, err := coreRebond(addr, height)
	ifExit(err)
	logger.Debugf("%v\n", tx)
	signAndBroadcast(c, chainID, nodeAddr, tx)
}
