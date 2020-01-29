package cmd

import (
	"os"

	"github.com/mainawycliffe/kamanda/configs"
	"github.com/mainawycliffe/kamanda/oauth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log kamanda into firebase",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet(configs.FirebaseRefreshTokenViperConfigKey) {
			email := viper.GetString(configs.FirebaseLoggedInUserEmailViperConfigKey)
			utils.StdOutSuccess("Already logged in as %s\n", email)
			os.Exit(1)
		}
		noLocalhostFlag, _ := cmd.Flags().GetBool("no-localhost")
		if !noLocalhostFlag {
			oauth.LoginWithLocalhost(false)
			os.Exit(0)
		}
		if err := oauth.LoginWithoutLocalhost(false); err != nil {
			utils.StdOutError("\n\n%s\n\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	// allow users to login without localhost
	loginCmd.Flags().Bool("no-localhost", false, " copy and paste a code instead of starting a local server for authentication")
}
