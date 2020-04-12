package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	shortened = false
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	output    = "json"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Version will output the current build information",
	Run: func(cmd *cobra.Command, args []string) {
		resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
		fmt.Print(resp)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&output, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")
}
