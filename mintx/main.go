package main

import (
	"fmt"
	"os"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/codegangsta/cli"
)

var (
	DefaultKeyDaemonHost = "localhost"
	DefaultKeyDaemonPort = "4767"
	DefaultKeyDaemonAddr = "http://" + DefaultKeyDaemonHost + ":" + DefaultKeyDaemonPort

	DefaultNodeRPCHost = "pinkpenguin.chaintest.net"
	DefaultNodeRPCPort = "46657"
	DefaultNodeRPCAddr = "http://" + DefaultNodeRPCHost + ":" + DefaultNodeRPCPort

	DefaultPubKey  string
	DefaultChainID string
)

// override the hardcoded defaults with env variables if they're set
func init() {
	signAddr := os.Getenv("MINTX_SIGN_ADDR")
	if signAddr != "" {
		DefaultKeyDaemonAddr = signAddr
	}

	nodeAddr := os.Getenv("MINTX_NODE_ADDR")
	if nodeAddr != "" {
		DefaultNodeRPCAddr = nodeAddr
	}

	pubKey := os.Getenv("MINTX_PUBKEY")
	if pubKey != "" {
		DefaultPubKey = pubKey
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
		signAddrFlag = cli.StringFlag{
			Name:  "sign-addr",
			Usage: "set the address of the eris-keys daemon",
			Value: DefaultKeyDaemonAddr,
		}

		nodeAddrFlag = cli.StringFlag{
			Name:  "node-addr",
			Usage: "set the address of the tendermint rpc server",
			Value: DefaultNodeRPCAddr,
		}

		pubkeyFlag = cli.StringFlag{
			Name:  "pubkey",
			Usage: "specify the pubkey",
			Value: DefaultPubKey,
		}

		chainidFlag = cli.StringFlag{
			Name:  "chainID",
			Usage: "specify the chainID",
			Value: DefaultChainID,
		}

		//----------------------------------------------------------------
		// optional action flags

		signFlag = cli.BoolFlag{
			Name:  "sign",
			Usage: "sign the transaction using the daemon at MINTX_SIGN_ADDR",
		}

		broadcastFlag = cli.BoolFlag{
			Name:  "broadcast",
			Usage: "broadcast the transaction using the daemon at MINTX_NODE_ADDR",
		}

		//----------------------------------------------------------------
		// tx data flags

		amtFlag = cli.StringFlag{
			Name:  "amt",
			Usage: "specify an amount",
		}

		nonceFlag = cli.StringFlag{
			Name:  "nonce",
			Usage: "set the account nonce",
		}

		addrFlag = cli.StringFlag{
			Name:  "addr",
			Usage: "specify an address",
		}

		nameFlag = cli.StringFlag{
			Name:  "name",
			Usage: "specify a name",
		}

		dataFlag = cli.StringFlag{
			Name:  "data",
			Usage: "specify some data",
		}

		dataFileFlag = cli.StringFlag{
			Name:  "data-file",
			Usage: "specify a file with some data",
		}

		toFlag = cli.StringFlag{
			Name:  "to",
			Usage: "specify an address to send to",
		}

		feeFlag = cli.StringFlag{
			Name:  "fee",
			Usage: "specify the fee to send",
		}

		gasFlag = cli.StringFlag{
			Name:  "gas",
			Usage: "specify the gas limit for a CallTx",
		}

		unbondtoFlag = cli.StringFlag{
			Name:  "unbond-to",
			Usage: "specify an address to unbond to",
		}

		heightFlag = cli.StringFlag{
			Name:  "height",
			Usage: "specify a height to unbond at",
		}

		//------------------------------------------------------------
		// main tx commands

		sendCmd = cli.Command{
			Name:   "send",
			Usage:  "mintx send --amt <amt> --to <addr>",
			Action: cliSend,
			Flags: []cli.Flag{
				signAddrFlag,
				nodeAddrFlag,

				chainidFlag,
				pubkeyFlag,
				addrFlag,

				signFlag,
				broadcastFlag,

				amtFlag,
				toFlag,
				nonceFlag,
			},
		}

		nameCmd = cli.Command{
			Name:   "name",
			Usage:  "mintx name --amt <amt> --name <name> --data <data>",
			Action: cliName,
			Flags: []cli.Flag{
				signAddrFlag,
				nodeAddrFlag,

				chainidFlag,
				pubkeyFlag,
				addrFlag,

				signFlag,
				broadcastFlag,

				amtFlag,
				nameFlag,
				dataFlag,
				dataFileFlag,
				feeFlag,
				nonceFlag,
			},
		}

		callCmd = cli.Command{
			Name:   "call",
			Usage:  "mintx call --amt <amt> --fee <fee> --gas <gas> --to <contract addr> --data <data>",
			Action: cliCall,
			Flags: []cli.Flag{
				signAddrFlag,
				nodeAddrFlag,

				chainidFlag,
				pubkeyFlag,
				addrFlag,

				signFlag,
				broadcastFlag,

				amtFlag,
				toFlag,
				dataFlag,
				feeFlag,
				gasFlag,
				nonceFlag,
			},
		}

		bondCmd = cli.Command{
			Name:   "bond",
			Usage:  "mintx bond --pubkey <pubkey> --amt <amt> --unbond-to <address>",
			Action: cliBond,
			Flags: []cli.Flag{
				signAddrFlag,
				nodeAddrFlag,

				chainidFlag,
				pubkeyFlag,
				addrFlag,

				signFlag,
				broadcastFlag,

				amtFlag,
				unbondtoFlag,
				nonceFlag,
			},
		}

		unbondCmd = cli.Command{
			Name:   "unbond",
			Usage:  "mintx unbond --addr <address> --height <block_height>",
			Action: cliUnbond,
			Flags: []cli.Flag{
				signAddrFlag,
				nodeAddrFlag,

				chainidFlag,

				signFlag,
				broadcastFlag,

				addrFlag,
				heightFlag,
			},
		}

		rebondCmd = cli.Command{
			Name:   "rebond",
			Usage:  "mintx rebond --addr <address> --height <block_height>",
			Action: cliRebond,
			Flags: []cli.Flag{
				signAddrFlag,
				nodeAddrFlag,

				chainidFlag,

				signFlag,
				broadcastFlag,

				addrFlag,
				heightFlag,
			},
		}

		/*
			inputCmd = cli.Command{
				Name:   "input",
				Usage:  "mintx input --pubkey <pubkey> --amt <amt> --nonce <nonce>",
				Action: cliInput,
				Flags: []cli.Flag{
					pubkeyFlag,
					amtFlag,
					nonceFlag,
					addrFlag,
				},
			}

			outputCmd = cli.Command{
				Name:   "output",
				Usage:  "mintx output --addr <addr> --amt <amt>",
				Action: cliOutput,
				Flags: []cli.Flag{
					addrFlag,
					amtFlag,
				},
			}*/

	)

	app := cli.NewApp()
	app.Name = "mintx"
	app.Usage = "Create and broadcast tendermint txs"
	app.Version = "0.0.1"
	app.Author = "Ethan Buchman"
	app.Email = "ethan@erisindustries.com"
	app.Commands = []cli.Command{
		// inputCmd,
		// outputCmd,
		sendCmd,
		nameCmd,
		callCmd,
		bondCmd,
		unbondCmd,
		rebondCmd,
		// dupeoutCmd,
	}
	app.Run(os.Args)

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
