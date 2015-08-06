package main

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/permission/types"
)

func cliStringsToInts(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) == 0 {
		Exit(fmt.Errorf("Please enter at least one `<permission>:<value>` pair like `send:0 call:1 create_account:1`"))
	}

	bp := types.ZeroBasePermissions

	for _, a := range args {
		spl := strings.Split(a, ":")
		if len(spl) != 2 {
			Exit(fmt.Errorf("arguments must be like `send:1`, not %s", a))
		}
		name, v := spl[0], spl[1]
		vi := v[0] - '0'
		pf, err := types.PermStringToFlag(name)
		IfExit(err)
		bp.Set(pf, vi > 0)
	}
	fmt.Println("Perms and SetBit (As Integers)")
	fmt.Printf("%d\t%d\n", bp.Perms, bp.SetBit)
	fmt.Println("\nPerms and SetBit (As Bitmasks)")
	fmt.Printf("%b\t%b\n", bp.Perms, bp.SetBit)
}

func coreIntsToStrings(perms, setbits types.PermFlag) map[string]bool {
	m := make(map[string]bool)

	for i := uint(0); i < types.NumPermissions; i++ {
		pf := types.PermFlag(1 << i)
		if pf&setbits > 0 {
			name := types.PermFlagToString(pf)
			m[name] = pf&perms > 0
		}
	}
	return m
}

func cliIntsToStrings(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) != 2 {
		Exit(fmt.Errorf("Please enter PermFlag and SetBit integers"))
	}

	pf, sb := args[0], args[1]
	perms, err := strconv.Atoi(pf)
	IfExit(err)
	setbits, err := strconv.Atoi(sb)
	IfExit(err)

	m := coreIntsToStrings(types.PermFlag(perms), types.PermFlag(setbits))
	for name, v := range m {
		fmt.Printf("%s: %v\n", name, v)
	}
}

func cliBBPB(cmd *cobra.Command, args []string) {
	pf := types.DefaultPermFlags
	fmt.Println("Perms and SetBit (As Integers)")
	fmt.Printf("%d\t%d\n", pf, pf)
	fmt.Println("\nPerms and SetBit (As Bitmasks)")
	fmt.Printf("%b\t%b\n", pf, pf)

	m := coreIntsToStrings(pf, pf)

	for name, v := range m {
		fmt.Printf("%s: %v\n", name, v)
	}
}

func cliAll(cmd *cobra.Command, args []string) {
	pf := types.AllPermFlags
	fmt.Println("Perms and SetBit (As Integers)")
	fmt.Printf("%d\t%d\n", pf, pf)
	fmt.Println("\nPerms and SetBit (As Bitmasks)")
	fmt.Printf("%b\t%b\n", pf, pf)

	m := coreIntsToStrings(pf, pf)

	for name, v := range m {
		fmt.Printf("%s: %v\n", name, v)
	}

}
