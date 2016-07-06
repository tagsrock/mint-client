package main

import (
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	moniker  string
	nodeAddr string
	seeds    string

	fast_sync bool
	skip_upnp bool

	db_backend string
	log_level  string
	rpcAddr    string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "mintconfig",
		Short: "a tool for generating config files",
		Long:  "use the flags to build a config file",
		Run:   cliConfig,
	}

	rootCmd.Flags().StringVarP(&moniker, "moniker", "", "golden_goose", "a moniker for your node; nice to have but not necessary")
	rootCmd.Flags().StringVarP(&nodeAddr, "p2p", "", "0.0.0.0:46656", "the p2p listening addr for your node")
	rootCmd.Flags().StringVarP(&rpcAddr, "rpc", "", "0.0.0.0:46657", "the RPC listening addr for your node")
	rootCmd.Flags().StringVarP(&seeds, "seeds", "", "", "seed address for new nodes; set to <ip:port> of the peer you'd like to connect to")

	rootCmd.Flags().BoolVarP(&fast_sync, "fast-sync", "", false, "catch up to an existing chain (true) or run the consensus protocol (false). WIP")
	rootCmd.Flags().BoolVarP(&skip_upnp, "skip-upnp", "", false, "skip UPNP port mapping")

	rootCmd.Flags().StringVarP(&db_backend, "db", "", "leveldb", "database back-end; options are: leveldb or memdb")
	rootCmd.Flags().StringVarP(&log_level, "log", "", "debug", "set the log level; options are: error < warn < notice < info < debug")

	rootCmd.Execute()
}
