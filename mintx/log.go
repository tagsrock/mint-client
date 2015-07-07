package main

import (
	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/log"
)

var logger *Logger

func init() {
	logger = AddLogger("mintx")
}
