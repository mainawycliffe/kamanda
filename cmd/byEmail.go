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
	Aliases: []string{"email"},
	Short:   "find a user by email address",
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError("atleast one email is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserEmailCriteria
		for _, email := range args {
			user, err := auth.GetUser(context.Background(), email, criteria)
			if err != nil {
				if firebase.IsUserNotFound(err) {
					utils.StdOutError("Not Found\t %s \t User was not found", email, err.Error())
					continue
				}
				utils.StdOutError("Error\t %s\t %s", email, err.Error())
				continue
			}
			//@todo something with the output
			utils.StdOutSuccess("Success\t%s\t%s", email, user.UID)
		}
		os.Exit(0)
	},
}

func init() {
	findUserCmd.AddCommand(byEmailCmd)
}
