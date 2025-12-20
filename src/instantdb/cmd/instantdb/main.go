package main

import (
	"fmt"
	"os"

	"github.com/db-toolkit/instant-db/src/instantdb/cmd/instantdb/commands"
)

var version = "0.1.0"

func main() {
	rootCmd := commands.GetRootCommand(version)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
