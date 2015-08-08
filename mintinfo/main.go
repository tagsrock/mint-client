package main

import (
	"fmt"
	"os"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	cclient "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core_client"
)

var (
	DefaultNodeRPCHost = "pinkpenguin.chaintest.net"
	DefaultNodeRPCPort = "46657"
	DefaultNodeRPCAddr = "http://" + DefaultNodeRPCHost + ":" + DefaultNodeRPCPort

	DefaultChainID string

	REQUEST_TYPE = "JSONRPC"
	client       cclient.Client
)

// override the hardcoded defaults with env variables if they're set
func init() {
	nodeAddr := os.Getenv("MINTX_NODE_ADDR")
	if nodeAddr != "" {
		DefaultNodeRPCAddr = nodeAddr
	}

	chainID := os.Getenv("MINTX_CHAINID")
	if chainID != "" {
		DefaultChainID = chainID
	}
}

func main() {

	// these are defined in here so we can update the
	// defaults with env variables first
	var (
		nodeAddrFlag string
		chainIDFlag  string
	)

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get a node's status",
		Run:   cliStatus,
	}

	var netInfoCmd = &cobra.Command{
		Use:   "net-info",
		Short: "Get a node's network info",
		Run:   cliNetInfo,
	}

	var genesisCmd = &cobra.Command{
		Use:   "genesis",
		Short: "Get a node's genesis.json",
		Run:   cliGenesis,
	}

	var validatorsCmd = &cobra.Command{
		Use:   "validators",
		Short: "List the chain's validator set",
		Run:   cliValidators,
	}

	var consensusCmd = &cobra.Command{
		Use:   "consensus",
		Short: "Dump a node's consensus state",
		Run:   cliConsensus,
	}

	var unconfirmedCmd = &cobra.Command{
		Use:   "unconfirmed",
		Short: "List the txs in a node's mempool",
		Run:   cliUnconfirmed,
	}

	var accountsCmd = &cobra.Command{
		Use:   "accounts",
		Short: "List all accounts on the chain, or specify an address",
		Run:   cliAccounts,
	}

	var namesCmd = &cobra.Command{
		Use:   "names",
		Short: "List all name reg entries on the chain",
		Run:   cliNames,
	}

	var blocksCmd = &cobra.Command{
		Use:   "blocks",
		Short: "Get a sequence of blocks between two heights, or get a single block by height",
		Run:   cliBlocks,
	}

	var storageCmd = &cobra.Command{
		Use:   "storage",
		Short: "Get the storage for an account, or for a particular key in that account's storage",
		Run:   cliStorage,
	}

	var callCmd = &cobra.Command{
		Use:   "call",
		Short: "Call an address with some data",
		Run:   cliCall,
	}

	var callCodeCmd = &cobra.Command{
		Use:   "call-code",
		Short: "Run some code on some data",
		Run:   cliCallCode,
	}

	var broadcastCmd = &cobra.Command{
		Use:   "broadcast",
		Short: "Broadcast some tx bytes",
		Run:   cliBroadcast,
	}

	var rootCmd = &cobra.Command{
		Use:   "mintinfo",
		Short: "Fetch data from a tendermint node via rpc",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			nodeAddr := nodeAddrFlag
			client = cclient.NewClient(nodeAddr, REQUEST_TYPE)
		},
	}

	// flags with env var defaults
	rootCmd.PersistentFlags().StringVarP(&nodeAddrFlag, "node-addr", "", DefaultNodeRPCAddr, "set the address of the tendermint rpc server")
	rootCmd.PersistentFlags().StringVarP(&chainIDFlag, "chainID", "", DefaultChainID, "specify the chainID")

	rootCmd.AddCommand(
		statusCmd,
		netInfoCmd,
		genesisCmd,
		validatorsCmd,
		consensusCmd,
		unconfirmedCmd,
		accountsCmd,
		namesCmd,
		blocksCmd,
		storageCmd,
		callCmd,
		callCodeCmd,
		broadcastCmd,
	)
	rootCmd.Execute()
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func ifExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
