package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	cclient "github.com/tendermint/tendermint/rpc/core_client"
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
		//----------------------------------------------------------------
		// flags with env var defaults
		nodeAddrFlag = cli.StringFlag{
			Name:  "node-addr",
			Usage: "set the address of the tendermint rpc server",
			Value: DefaultNodeRPCAddr,
		}

		chainidFlag = cli.StringFlag{
			Name:  "chainID",
			Usage: "specify the chainID",
			Value: DefaultChainID,
		}

		//----------------------------------------------------------------

		statusCmd = cli.Command{
			Name:   "status",
			Usage:  "Get a node's status",
			Action: cliStatus,
		}

		netInfoCmd = cli.Command{
			Name:   "net-info",
			Usage:  "Get a node's network info",
			Action: cliNetInfo,
		}

		genesisCmd = cli.Command{
			Name:   "genesis",
			Usage:  "Get a node's genesis.json",
			Action: cliGenesis,
		}

		validatorsCmd = cli.Command{
			Name:   "validators",
			Usage:  "List the chain's validator set",
			Action: cliValidators,
		}

		consensusCmd = cli.Command{
			Name:   "consensus",
			Usage:  "Dump a node's consensus state",
			Action: cliConsensus,
		}

		unconfirmedCmd = cli.Command{
			Name:   "unconfirmed",
			Usage:  "List the txs in a node's mempool",
			Action: cliUnconfirmed,
		}

		accountsCmd = cli.Command{
			Name:   "accounts",
			Usage:  "List all accounts on the chain, or specify an address",
			Action: cliAccounts,
		}

		namesCmd = cli.Command{
			Name:   "names",
			Usage:  "List all name reg entries on the chain",
			Action: cliNames,
		}

		blocksCmd = cli.Command{
			Name:   "blocks",
			Usage:  "Get a sequence of blocks between two heights, or get a single block by height",
			Action: cliBlocks,
		}

		storageCmd = cli.Command{
			Name:   "storage",
			Usage:  "Get the storage for an account, or for a particular key in that account's storage",
			Action: cliStorage,
		}

		callCmd = cli.Command{
			Name:   "call",
			Usage:  "Call an address with some data",
			Action: cliCall,
		}

		callCodeCmd = cli.Command{
			Name:   "call-code",
			Usage:  "Run some code on some data",
			Action: cliCallCode,
		}

		broadcastCmd = cli.Command{
			Name:   "broadcast",
			Usage:  "Broadcast some tx bytes",
			Action: cliBroadcast,
		}
	)

	app := cli.NewApp()
	app.Name = "mintinfo"
	app.Usage = "Fetch data from a tendermint node via rpc"
	app.Version = "0.0.1"
	app.Author = "Ethan Buchman"
	app.Email = "ethan@erisindustries.com"
	app.Commands = []cli.Command{
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
	}
	app.Flags = []cli.Flag{
		nodeAddrFlag,
		chainidFlag,
	}
	app.Before = before

	app.Run(os.Args)
}

func before(c *cli.Context) error {
	nodeAddr := c.String("node-addr")
	client = cclient.NewClient(nodeAddr, REQUEST_TYPE)
	return nil
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
