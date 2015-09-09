package utils

import (
	cfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
	tmcfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config/tendermint"
)

func init() {
	cfg.ApplyConfig(tmcfg.GetConfig(""))
}
