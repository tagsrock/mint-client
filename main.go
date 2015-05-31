package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	cfg "github.com/tendermint/tendermint/config"
	tmcfg "github.com/tendermint/tendermint/config/tendermint"
)

var (
	DefaultKeyDaemonHost = "localhost"
	DefaultKeyDaemonPort = "4767"
	DefaultKeyDaemonAddr = "http://" + DefaultKeyDaemonHost + ":" + DefaultKeyDaemonPort
)

func main() {
	app := cli.NewApp()
	app.Name = "mintx"
	app.Usage = "Create and broadcast tendermint txs"
	app.Version = "0.0.1"
	app.Author = "Ethan Buchman"
	app.Email = "ethan@erisindustries.com"
	app.Before = before
	app.Commands = []cli.Command{
		inputCmd,
		outputCmd,
		sendCmd,
		nameCmd,
		// callCmd,
		// bondCmd,
		// unbondCmd,
		// rebondCmd,
		// dupeoutCmd,
	}

	app.Run(os.Args)

}

var (
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
	}

	sendCmd = cli.Command{
		Name:   "send",
		Usage:  "mintx send --pubkey <pubkey> --amt <amt> --nonce <nonce> --to <addr>",
		Action: cliSend,
		Flags: []cli.Flag{
			pubkeyFlag,
			amtFlag,
			nonceFlag,
			addrFlag,
			toFlag,
			chainidFlag,
			signFlag,
			broadcastFlag,
		},
	}
	/*	sendCmd = cli.Command{
		Name:   "send",
		Usage:  "mintx send --inputs <input1,input2,...> --outputs <output1,output2,...>",
		Action: cliSend,
		Flags: []cli.Flag{
			inputsFlag,
			outputsFlag,
		},
	}*/

	nameCmd = cli.Command{
		Name:   "name",
		Usage:  "mintx name --pubkey <pubkey> --amt <amt> --nonce <nonce> --name <name> --data <data>",
		Action: cliName,
		Flags: []cli.Flag{
			pubkeyFlag,
			amtFlag,
			nonceFlag,
			addrFlag,
			nameFlag,
			dataFlag,
			feeFlag,
			chainidFlag,
			signFlag,
			broadcastFlag,
		},
	}

	signFlag = cli.StringFlag{
		Name:  "sign",
		Usage: "specify the rpc address of the signing daemon for signing the tx",
		Value: DefaultKeyDaemonAddr,
	}

	broadcastFlag = cli.StringFlag{
		Name:  "broadcast",
		Usage: "specify the rpc address of a node for broadcasting the tx",
	}

	pubkeyFlag = cli.StringFlag{
		Name:  "pubkey",
		Usage: "specify the pubkey",
	}

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

	toFlag = cli.StringFlag{
		Name:  "to",
		Usage: "specify an address to send to",
	}

	feeFlag = cli.StringFlag{
		Name:  "fee",
		Usage: "specify the fee to send",
	}

	inputsFlag = cli.StringFlag{
		Name:  "inputs",
		Usage: "csv list of hex encoded inputs",
	}

	outputsFlag = cli.StringFlag{
		Name:  "outputs",
		Usage: "csv list of hex encoded outputs",
	}

	chainidFlag = cli.StringFlag{
		Name:  "chainID",
		Usage: "specify the chainID",
	}
)

func before(c *cli.Context) error {
	config := tmcfg.GetConfig("")
	cfg.ApplyConfig(config) // Notify modules of new config
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
