package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
		ifExit(addTinyDNSARecord(entry.Name, entry.Data))
	}

	// done adding entries. commit them
	ifExit(makeTinyDNSRecords())
}

func makeTinyDNSRecords() error {
	cmd := exec.Command("make")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func addTinyDNSARecord(fqdn, addr string) error {
	fmt.Println("Running add host", fqdn, addr, "...")
	cmd := exec.Command("./add-host", fqdn, addr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("\t ... running add-alias")
		cmd := exec.Command("./add-alias", fqdn, addr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

func cliRun(cmd *cobra.Command, args []string) {
	// parse the tinydns data
	dnsData, err := TinyDNSDataFromFile(tinydnsDataFileFlag)
	ifExit(err)

	fetchAndUpdateRecords(dnsData)

	ticker := time.Tick(time.Second * time.Duration(updateEveryFlag))
	for {
		select {
		case <-ticker:
			fetchAndUpdateRecords(dnsData)
		}
	}
}

func fetchAndUpdateRecords(dnsData TinyDNSData) {
	// get all dns entries from chain
	dnsEntries, err := getDNSEntries()
	if err != nil {
		fmt.Println("Error getting dns entries", err)
	}

	anyUpdates := false
	for _, entry := range dnsEntries {
		record, ok := dnsData[entry.Name]

		toUpdate := true
		// if we have it and nothings changed, don't update
		if ok && record.FQDN == entry.Name && record.Address == entry.Data {
			toUpdate = false
		}

		if toUpdate {
			anyUpdates = true
			addTinyDNSARecord(entry.Name, entry.Data)
			// TODO: handle record types better
			dnsData[entry.Name] = &ResourceRecord{"", entry.Name, entry.Data}
		}
	}

	if anyUpdates {
		// done adding entries. commit them
		if err = makeTinyDNSRecords(); err != nil {
			fmt.Println("Error rebuilding data.cdb", err)
		}
	} else {
		fmt.Println("No new updates")
	}

}
