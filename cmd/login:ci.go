package cmd

import (
	"os"

	"github.com/mainawycliffe/kamanda/oauth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// loginCICmd represents the login:ci command
var loginCICmd = &cobra.Command{
	Use:   "login:ci",
	Short: "generate an access token for use in non-interactive environments",
	Run: func(cmd *cobra.Command, args []string) {
		noLocalhostFlag, _ := cmd.Flags().GetBool("no-localhost")
		if !noLocalhostFlag {
			oauth.LoginWithLocalhost(true)
			os.Exit(0)
		}
		if err := oauth.LoginWithoutLocalhost(true); err != nil {
			utils.StdOutError(os.Stderr, "\n\n%s\n\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(loginCICmd)
	// allow users to login without localhost
	loginCICmd.Flags().Bool("no-localhost", false, "copy and paste a code instead of starting a local server for authentication")
}
