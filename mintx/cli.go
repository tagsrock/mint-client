package main

import (
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/mintx/core"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

func cliSend(cmd *cobra.Command, args []string) {
	tx, err := core.Send(nodeAddrFlag, signAddrFlag, pubkeyFlag, addrFlag, toFlag, amtFlag, nonceFlag)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func cliName(cmd *cobra.Command, args []string) {
	if dataFlag != "" && dataFileFlag != "" {
		common.Exit(fmt.Errorf("Please specify only one of --data and --data-file"))
	}
	data := dataFlag
	if dataFlag == "" && dataFileFlag != "" {
		b, err := ioutil.ReadFile(dataFileFlag)
		common.IfExit(err)
		data = string(b)
	}
	tx, err := core.Name(nodeAddrFlag, signAddrFlag, pubkeyFlag, addrFlag, amtFlag, nonceFlag, feeFlag, nameFlag, data)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func cliCall(cmd *cobra.Command, args []string) {
	tx, err := core.Call(nodeAddrFlag, signAddrFlag, pubkeyFlag, addrFlag, toFlag, amtFlag, nonceFlag, gasFlag, feeFlag, dataFlag)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func cliPermissions(cmd *cobra.Command, args []string) {
	// all functions take at least 2 args (+ name)
	if len(args) < 3 {
		common.Exit(fmt.Errorf("Please enter the permission function you'd like to call, followed by it's arguments"))
	}
	permFunc := args[0]
	tx, err := core.Permissions(nodeAddrFlag, signAddrFlag, pubkeyFlag, addrFlag, nonceFlag, permFunc, args[1:])
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func cliBond(cmd *cobra.Command, args []string) {
	tx, err := core.Bond(nodeAddrFlag, signAddrFlag, pubkeyFlag, unbondtoFlag, amtFlag, nonceFlag)
	common.IfExit(err)

	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func cliUnbond(cmd *cobra.Command, args []string) {
	tx, err := core.Unbond(addrFlag, heightFlag)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func cliRebond(cmd *cobra.Command, args []string) {
	tx, err := core.Rebond(addrFlag, heightFlag)
	common.IfExit(err)
	logger.Debugf("%v\n", tx)
	unpackSignAndBroadcast(core.SignAndBroadcast(chainidFlag, nodeAddrFlag, signAddrFlag, tx, signFlag, broadcastFlag, waitFlag))
}

func unpackSignAndBroadcast(result *core.TxResult, err error) {
	common.IfExit(err)
	if result == nil {
		// if we don't provide --sign or --broadcast
		return
	}
	fmt.Printf("Transaction Hash: %X\n", result.Hash)
	if result.Address != nil {
		fmt.Printf("Contract Address: %X\n", result.Address)
	}
	if result.Return != nil {
		fmt.Printf("Block Hash: %X\n", result.BlockHash)
		fmt.Printf("Return Value: %X\n", result.Return)
		fmt.Printf("Exception: %s\n", result.Exception)
	}
}
