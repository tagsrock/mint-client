package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	cfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
	tmcfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config/tendermint"
)

var (
	usr, _ = user.Current()

	DefaultDir = path.Join(usr.HomeDir, ".tendermint")
	config     = tmcfg.GetConfig("")
)

func init() {
	cfg.ApplyConfig(config)
}

func main() {
	var dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump tendermint state to json files",
		Run:   cliDump,
	}
	var restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore tendermint state from json files",
		Long:  "mintdump restore <path/to/file> <new chainID>",
		Run:   cliRestore,
	}

	var rootCmd = &cobra.Command{Use: "mintdump"}
	rootCmd.AddCommand(dumpCmd, restoreCmd)
	rootCmd.Execute()
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func ifExit(err error) {
	if err != nil {
		exit(err)
	}
}
