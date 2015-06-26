package main

import (
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	NameFlag string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "genesis",
		Short: "Create a set of keys and a genesis file from them",
		Long:  "Create a set of keys and a genesis file from them",
		Run:   cliGenesis,
	}
	rootCmd.Flags().StringVarP(&NameFlag, "name", "n", "tendermint_test", "name for the chain (chain id)")
	rootCmd.Execute()
}
