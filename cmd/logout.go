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

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout kamanda from Firebase",
	Run: func(cmd *cobra.Command, args []string) {
		email := viper.GetString(configs.FirebaseLoggedInUserEmailViperConfigKey)
		if !viper.IsSet(configs.FirebaseRefreshTokenViperConfigKey) {
			fmt.Fprint(os.Stdout, aurora.Sprintf("%s\n", aurora.Red("You are not logged in!")))
			os.Exit(1)
		}
		if err := oauth.RevokeRefreshToken(); err != nil {
			fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Red("%s\n"), err.Error()))
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Green("Logged out from %s\n\n"), email))
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
