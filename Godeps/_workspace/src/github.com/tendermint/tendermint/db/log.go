package db

import (
	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
)

var logga *Logger

func init() {
	logga = AddLogger("mintdb")
}
