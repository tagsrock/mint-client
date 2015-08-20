package main

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

func setConfig(cmd *cobra.Command, args []string) {

	checkFlags(nodeAddr, seeds, db_backend, log_level, rpcAddr, fast_sync)

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

func checkFlags(node, seeds, db, log, rpc string, sync bool) {

	if node != "" {
		p := strings.Split(node, ":")
		if len(p) != 2 {
			Exit(fmt.Errorf("--p2p requires a port"))
		}
		_, err := strconv.Atoi(p[1])
		if err != nil {
			Exit(fmt.Errorf("specified port must be number"))
		}
	}

	if rpc != "" {
		r := strings.Split(rpc, ":")
		if len(r) != 2 {
			Exit(fmt.Errorf("--rpc requires a port"))
		}
		_, err := strconv.Atoi(r[1])
		if err != nil {
			Exit(fmt.Errorf("specified port must be number"))
		}
	}

	if seeds != "" {
		s := strings.Split(seeds, ":")
		if len(s) != 2 {
			Exit(fmt.Errorf("--seeds requires a port"))
		}
		_, err := strconv.Atoi(s[1])
		if err != nil {
			Exit(fmt.Errorf("specified port must be number"))
		}
	}

	if sync != true && sync != false {
		Exit(fmt.Errorf("--fast-sync must be true or false"))
	}

	if db != "memdb" && db != "leveldb" {
		Exit(fmt.Errorf("--db must be either leveldb or memdb"))
	}

	if log != "error" && log != "warn" && log != "notice" && log != "info" && log != "debug" {
		Exit(fmt.Errorf("--log must be one of: error, warn, notice, info, debug"))
	}

}
