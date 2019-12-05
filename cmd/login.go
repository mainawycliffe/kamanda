package cmd

import (
	"github.com/mainawycliffe/kamanda/oauth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		noLocalhost, _ := cmd.Flags().GetBool("no-localhost")
		if noLocalhost {
			oauth.LoginWithoutLocalhost()
		} else {
			oauth.LoginWithLocalhost()
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// allow users to login without localhost
	loginCmd.Flags().Bool("no-localhost", false, " copy and paste a code instead of starting a local server for authentication")
}
