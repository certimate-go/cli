package main

import (
	"os"

	"github.com/certimate-go/cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.BuildDate = date

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
