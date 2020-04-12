package main

import (
	"fmt"

	"github.com/mainawycliffe/kamanda/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "goreleaser"
)

func main() {
	fmt.Printf("%s %s %s %s", version, commit, date, builtBy)
	cmd.Execute()
}
