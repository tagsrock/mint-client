package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/permission/types"
)

func cliStringsToInts(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) == 0 {
		exit(fmt.Errorf("Please enter at least one `<permission>:<value>` pair like `send:0 call:1 create_account:1`"))
	}

	bp := types.NewBasePermissions()

	for _, a := range args {
		spl := strings.Split(a, ":")
		if len(spl) != 2 {
			exit(fmt.Errorf("arguments must be like `send:1`, not %s", a))
		}
		name, v := spl[0], spl[1]
		vi := v[0] - '0'
		pf, err := types.PermStringToFlag(name)
		ifExit(err)
		bp.Set(pf, vi > 0)
	}
	fmt.Println("Perms and SetBit (As Integers)")
	fmt.Printf("%d\t%d\n", bp.Perms, bp.SetBit)
	fmt.Println("\nPerms and SetBit (As Bitmasks)")
	fmt.Printf("%b\t%b\n", bp.Perms, bp.SetBit)
}

func cliIntsToStrings(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) != 2 {
		exit(fmt.Errorf("Please enter PermFlag and SetBit integers"))
	}

	pf, sb := args[0], args[1]
	perms, err := strconv.Atoi(pf)
	ifExit(err)
	setbits, err := strconv.Atoi(sb)
	ifExit(err)

	m := make(map[string]bool)

	for i := uint(0); i < types.NumBasePermissions; i++ {
		pf := 1 << i
		if pf&setbits > 0 {
			name, _ := types.PermFlagToString(types.PermFlag(pf))
			m[name] = pf&perms > 0
		}
	}
	for name, v := range m {
		fmt.Printf("%s: %v\n", name, v)
	}

}
