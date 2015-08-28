package main

import (
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/mintx/core"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

// do we really need these?
/*
func cliInput(c *cli.Context) {
	pubkey, amtS, nonceS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("addr")
	input, err := coreInput(pubkey, amtS, nonceS, addr)
	common.IfExit(err)
	fmt.Printf("%s\n", input)
}

func cliOutput(c *cli.Context) {
	addr, amtS := c.String("addr"), c.String("amt")
	output, err := coreOutput(addr, amtS)
	common.IfExit(err)
	fmt.Printf("%s\n", output)
}
*/

func cliSend(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	pubkey, amtS, nonceS, addr, toAddr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("addr"), c.String("to")
	tx, err := core.Send(nodeAddr, pubkey, addr, toAddr, amtS, nonceS)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func cliName(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	fmt.Println("CHAIN ID FROM NAME:", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")

	if c.IsSet("data") && c.IsSet("data-file") {
		common.Exit(fmt.Errorf("Please specify only one of --data and --data-file"))
	}
	name, data, dataFile := c.String("name"), c.String("data"), c.String("data-file")
	if data == "" && dataFile != "" {
		b, err := ioutil.ReadFile(dataFile)
		common.IfExit(err)
		data = string(b)
	}
	tx, err := core.Name(nodeAddr, pubkey, addr, amtS, nonceS, feeS, name, data)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func cliCall(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	pubkey, amtS, nonceS, feeS, addr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("fee"), c.String("addr")
	toAddr, gasS, data := c.String("to"), c.String("gas"), c.String("data")
	tx, err := core.Call(nodeAddr, pubkey, addr, toAddr, amtS, nonceS, gasS, feeS, data)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func cliPermissions(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	pubkey, nonceS, addr := c.String("pubkey"), c.String("nonce"), c.String("addr")

	// all functions take at least 2 args (+ name)
	if len(c.Args()) < 3 {
		common.Exit(fmt.Errorf("Please enter the permission function you'd like to call, followed by it's arguments"))
	}
	permFunc := c.Args()[0]
	tx, err := core.Permissions(nodeAddr, pubkey, addr, nonceS, permFunc, c.Args()[1:])
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func cliNewAccount(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	/*
		chainID, nodeAddr := c.String("chainID"), c.String("node-addr")
		pubkey := c.String("pubkey")

		tx, err := coreNewAccount(nodeAddr,signAddr, pubkey, chainID)
		common.IfExit(err)

		logger.Debugf("%v\n", tx)
		unpackSignAndBroadcast(core.SignAndBroadcast( chainID, nodeAddr,signAddr, tx, sign, broadcast, wait)
	*/
}

func cliBond(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	pubkey, amtS, nonceS, unbondAddr := c.String("pubkey"), c.String("amt"), c.String("nonce"), c.String("unbond-to")
	tx, err := core.Bond(nodeAddr, pubkey, unbondAddr, amtS, nonceS)
	common.IfExit(err)

	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func cliUnbond(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	addr, height := c.String("addr"), c.String("height")
	tx, err := core.Unbond(addr, height)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func cliRebond(c *cli.Context) {
	config.Set("chain_id", c.String("chainID"))
	chainID, nodeAddr, signAddr := c.String("chainID"), c.String("node-addr"), c.String("sign-addr")
	sign, broadcast, wait := c.Bool("sign"), c.Bool("broadcast"), c.Bool("wait")
	addr, height := c.String("addr"), c.String("height")
	tx, err := core.Rebond(addr, height)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainID, nodeAddr, signAddr, tx, sign, broadcast, wait))
}

func unpackSignAndBroadcast(result *core.TxResult, err error) {
	common.IfExit(err)
	if result == nil {
		// if we don't provide --sign or --broadcast
		return
	}
	fmt.Printf("Transaction Hash: %X\n", result.Hash)
	if result.Return != nil {
		fmt.Printf("Block Hash: %X\n", result.BlockHash)
		fmt.Printf("Return Value: %X\n", result.Return)
		fmt.Printf("Exception: %s\n", result.Exception)
	}
}
