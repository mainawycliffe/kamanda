package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mainawycliffe/kamanda/configs"
	"github.com/mainawycliffe/kamanda/oauth"
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
			fmt.Fprint(os.Stdout, aurora.Sprintf("Already logged in as %s\n", aurora.Green(email)))
			os.Exit(1)
		}
		noLocalhostFlag, _ := cmd.Flags().GetBool("no-localhost")
		if !noLocalhostFlag {
			oauth.LoginWithLocalhost(false)
			os.Exit(0)
		}
		if err := oauth.LoginWithoutLocalhost(false); err != nil {
			fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Red("\n\n%s\n\n"), err.Error()))
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
