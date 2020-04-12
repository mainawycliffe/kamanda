package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/cheynewallace/tabby"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "goreleaser"
	output  = "text"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Version will output the current build information",
	Run: func(cmd *cobra.Command, args []string) {
		kamandaVersion := map[string]string{
			"version":    version,
			"commitHash": commit,
			"builtBy":    builtBy,
			"built":      date,
			"OS/Arch":    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		}
		if output == "text" {
			t := tabby.New()
			t.AddLine("Version:", version)
			t.AddLine("Release Date:", date)
			t.AddLine("Commit Hash:", commit)
			t.AddLine("Built by:", builtBy)
			t.AddLine("OS/Arch:", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
			t.Print()
			os.Exit(0)
		}
		if output == "json" {
			output, err := json.Marshal(kamandaVersion)
			if err != nil {
				utils.StdOutError(os.Stderr, "Error marsalling json %s", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", output)
			os.Exit(0)
		}
		if output == "yaml" {
			output, err := yaml.Marshal(kamandaVersion)
			if err != nil {
				utils.StdOutError(os.Stderr, "Error marsalling yaml %s", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", output)
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of 'Text', 'yaml' or 'json'.")
}
