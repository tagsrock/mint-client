package main

import (
	"fmt"
	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

func setConfig(cmd *cobra.Command, args []string) {

	checkFlags(db_backend, log_level)

	//XXX this formating necessary for properly formatted stdout
	var Config = fmt.Sprintf(
		`# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

moniker = "%s"
node_laddr = "%s"
seeds = "%s"
fast_sync = %t
db_backend = "%s"
log_level = "%s"
rpc_laddr = "%s"`, moniker, nodeAddr, seeds, fast_sync, db_backend, log_level, rpcAddr)

	fmt.Printf("%v\n", Config)

}

func checkFlags(db, log string) {

	if db != "memdb" && db != "leveldb" {
		Exit(fmt.Errorf("--db must be either leveldb or memdb"))
	}

	if log != "error" && log != "warn" && log != "notice" && log != "info" && log != "debug" {
		Exit(fmt.Errorf("--log must be one of: error, warn, notice, info, debug"))
	}

}
