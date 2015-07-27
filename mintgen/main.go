package main

import (
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	DirFlag    string
	PubkeyFlag string
	//AddrsFlag  string
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
	randomCmd.Flags().StringVarP(&DirFlag, "dir", "d", "", "Directory to save genesis and priv_validators in. Default is ~/.eris/data/<chain_id>")

	//XXX uses pubkey until I figure out how to do conversion
	var multiCmd = &cobra.Command{
		Use:   "multi",
		Short: "mintgen multi <chain_id> --pub <pub_1> <pub_2> <pub_N>",
		Long:  "Create a genesis.json with <chain_id> and N <pub>'s passed in, seperated by a space; --pub is req'd",
		Run:   cliMulti,
	}
	multiCmd.Flags().StringVarP(&PubkeyFlag, "pub", "", "", "pubkeys to include when generating genesis.json. flag is req'd")
	multiCmd.Flags().StringVarP(&DirFlag, "dir", "d", "", "Directory to save genesis.json in. Default is ~/.eris/data/<chain_id>")

	var rootCmd = &cobra.Command{
		Use:   "mintgen",
		Short: "a tool for generating tendermint genesis files",
		Long:  "a tool for generating tendermint genesis files",
	}
	rootCmd.AddCommand(singleCmd, randomCmd, multiCmd)
	rootCmd.Execute()
}
