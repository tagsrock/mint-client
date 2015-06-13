package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path"
)

var (
	usr, _ = user.Current()

	DefaultKeyStore = path.Join(usr.HomeDir, ".decerver", "keys")
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mintkey",
		Short: "Convert an eris-keys key to a priv_validator.json",
		Run:   cliConvertAddressToPrivValidator,
	}
	rootCmd.Execute()
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func ifExit(err error) {
	if err != nil {
		exit(err)
	}
}
