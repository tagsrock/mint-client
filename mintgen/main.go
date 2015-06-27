package main

import (
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	SingleFlag bool
	NameFlag   string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "mintgen",
		Short: "Create a set of keys and a genesis file from them",
		Long:  "Create a set of keys and a genesis file from them",
		Run:   cliGenesis,
	}
	rootCmd.Flags().StringVarP(&NameFlag, "name", "n", "tendermint_test", "name for the chain (chain id)")
	rootCmd.Flags().BoolVarP(&SingleFlag, "single", "s", false, "create a genesis.json with a single key")
	rootCmd.Execute()
}
