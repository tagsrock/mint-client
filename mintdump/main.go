package main

import (
	"fmt"
	"os"
	"path"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	cfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
	tmcfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config/tendermint"
)

var (
	config = tmcfg.GetConfig("")

	DataDirFlag string
)

func init() {
	cfg.ApplyConfig(config)
}

func main() {

	var dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "Dump tendermint state to json files",
		Long:  "mintdump dump > [json file]",
		Run:   cliDump,
	}
	var restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore tendermint state from json files",
		Long:  "mintdump restore [new chainID] < [path/to/file] ",
		Run:   cliRestore,
	}
	dumpCmd.Flags().StringVarP(&DataDirFlag, "data-dir", "d", "", "Path to tendermint data directory")
	restoreCmd.Flags().StringVarP(&DataDirFlag, "data-dir", "d", "", "Path to tendermint data directory")

	var rootCmd = &cobra.Command{Use: "mintdump"}
	rootCmd.PersistentPreRun = before
	rootCmd.AddCommand(dumpCmd, restoreCmd)
	rootCmd.Execute()
}

func before(cmd *cobra.Command, args []string) {
	if DataDirFlag != "" {
		if _, err := os.Stat(path.Join(DataDirFlag, "state.db")); err != nil {
			exit(fmt.Errorf("Could not find state.db folder in %s", DataDirFlag))
		}
		config.Set("db_dir", DataDirFlag)
		cfg.ApplyConfig(config)
	}
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
