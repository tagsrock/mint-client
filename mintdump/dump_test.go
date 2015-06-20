package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"testing"

	cfg "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
	dbm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/db"
	ptypes "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/permission/types"
	sm "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/state"
)

var TestDir = path.Join(os.Getenv("GOPATH"), "src", "github.com", "eris-ltd", "mint-client", "mintdump", "test")

/*
// XXX: Keeping leveldb files in version control is not nice ....
func TestDumpRestore(t *testing.T) {
	// load a test dir
	config.Set("db_dir", path.Join(TestDir, "data1"))
	cfg.ApplyConfig(config) // Notify modules of new config
	stateDB := dbm.GetDB("state")
	st := sm.LoadState(stateDB)
	stHash := st.Hash()
	dump := CoreDump()

	// restore to a memdir
	config.Set("db_backend", "memdb")
	cfg.ApplyConfig(config) // Notify modules of new config
	CoreRestore("", dump)

	stateDB = dbm.GetDB("state")
	st = sm.LoadState(stateDB)
	if bytes.Compare(stHash, st.Hash()) != 0 {
		t.Fatalf("State hash mismatch. Got %X, expected %X", st.Hash(), stHash)
	}
}
*/

func TestRestoreDump(t *testing.T) {
	b, err := ioutil.ReadFile(path.Join(TestDir, "data1.json"))
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.Trim(b, "\n")

	// restore to a memdir
	config.Set("db_backend", "memdb")
	cfg.ApplyConfig(config) // Notify modules of new config
	CoreRestore("", b)

	stateDB := dbm.GetDB("state")
	st := sm.LoadState(stateDB)
	acc := st.GetAccount(ptypes.GlobalPermissionsAddress)
	fmt.Println(acc)

	dump := CoreDump()

	if bytes.Compare(b, dump) != 0 {
		ld, lb := len(dump), len(b)
		max := int(math.Max(float64(ld), float64(lb)))
		n := 100
		for i := 0; i < max/n; i++ {
			dd := dump[i*n : (i+1)*n]
			bb := b[i*n : (i+1)*n]
			if bytes.Compare(dd, bb) != 0 {
				t.Fatalf("Error in dumps! Got \n\n\n\n %s \n\n\n\n Expected \n\n\n\n %s", dd, bb)
			}
		}
	}
}
