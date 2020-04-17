package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// byEmailCmd represents the byEmail command
var byEmailCmd = &cobra.Command{
	Use:     "byEmail",
	Aliases: []string{"email", "by-email"},
	Short:   "Find a Firebase Auth user by email address",
	Example: `kamanda auth find by-email email@example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of UIDs
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one email is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserEmailCriteria
		for _, email := range args {
			user, err := auth.GetUser(context.Background(), email, criteria)
			if err != nil {
				if firebase.IsUserNotFound(err) {
					utils.StdOutError(os.Stderr, "Not Found\t %s \t User was not found", email, err.Error())
					continue
				}
				utils.StdOutError(os.Stderr, "Error\t %s\t %s", email, err.Error())
				continue
			}
			//@todo something with the output
			utils.StdOutSuccess(os.Stdout, "Success\t%s\t%s", email, user.UID)
		}
		os.Exit(0)
	},
}

func init() {
	findUserCmd.AddCommand(byEmailCmd)
}
