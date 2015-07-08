package main

import (
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	DirFlag    string
	PubkeyFlag string
)

func main() {
	var singleCmd = &cobra.Command{
		Use:   "single",
		Short: "mintgen single <chain_id>",
		Long:  "Create a genesis.json with <chain_id> from a priv_validator.json passed on stdin",
		Run:   cliSingle,
	}
	singleCmd.Flags().StringVarP(&PubkeyFlag, "pub", "p", "", "pubkey to use instead of a priv_validator.json")

	var randomCmd = &cobra.Command{
		Use:   "random",
		Short: "mintgen random [flags] <N> <chain_id>",
		Long:  "Create <N> keys and a genesis.json with corresponding validators and chain_id <name>",
		Run:   cliRandom,
	}
	randomCmd.Flags().StringVarP(&DirFlag, "dir", "d", "", "Directory to save genesis and priv_validators in")

	var rootCmd = &cobra.Command{
		Use:   "mintgen",
		Short: "a tool for generating tendermint genesis files",
		Long:  "a tool for generating tendermint genesis files",
	}
	rootCmd.AddCommand(singleCmd, randomCmd)
	rootCmd.Execute()
}
