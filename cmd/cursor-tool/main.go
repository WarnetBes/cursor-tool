package main

import (
	"fmt"
	"os"

	"github.com/WarnetBes/cursor-tool/cmd/cursor-tool/commands"
)

var Version = "dev"

func main() {
	commands.Version = Version
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
