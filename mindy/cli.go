package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
)

func validateDNSEntrySimple(entry *types.NameRegEntry) error {
	spl := strings.Split(entry.Name, ".")
	if len(spl) < 3 {
		return fmt.Errorf("A valid name must have at least a subdomain, host name, and tld")
	}
	spl = strings.Split(entry.Data, ".")
	if len(spl) != 4 {
		return fmt.Errorf("Data must be a valid ipv4 address")
	}
	return nil
}

func getDNSEntries() ([]*types.NameRegEntry, error) {
	r, err := client.ListNames()
	if err != nil {
		return nil, err
	}
	dnsEntries := []*types.NameRegEntry{}
	for _, entry := range r.Names {
		if err := validateDNSEntrySimple(entry); err == nil {
			dnsEntries = append(dnsEntries, entry)
		}
	}
	return dnsEntries, nil
}

func cliListNames(cmd *cobra.Command, args []string) {
	dnsEntries, err := getDNSEntries()
	ifExit(err)
	s, err := formatOutput(args, 1, dnsEntries)
	ifExit(err)
	fmt.Println(s)
}

func cliCatchup(cobraCmd *cobra.Command, args []string) {
	dnsEntries, err := getDNSEntries()
	ifExit(err)

	err = os.Chdir(DefaultTinyDNSDir)
	ifExit(err)

	// KISS
	// for each entry, try add-host
	// if fails, run add-alias
	// TODO: parse/update the data file manually?
	for _, entry := range dnsEntries {
		fmt.Println("Running add host", entry.Name, entry.Data, "...")
		cmd := exec.Command("./add-host", entry.Name, entry.Data)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("\t ... running add-alias")
			cmd := exec.Command("./add-alias", entry.Name, entry.Data)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			ifExit(err)
		}
	}

	// done adding entries. commit the,
	cmd := exec.Command("make")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	ifExit(err)
}

func cliRun(cmd *cobra.Command, args []string) {

}
