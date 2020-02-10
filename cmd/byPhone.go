package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// byPhoneCmd represents the byPhone command
var byPhoneCmd = &cobra.Command{
	Use:     "byPhone",
	Aliases: []string{"phone"},
	Short:   "find a user by phone number",
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of phone numbers
		if len(args) == 0 {
			utils.StdOutError("at least one Firebase user UID is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserUIDCriteria
		for _, phone := range args {
			user, err := auth.GetUser(context.Background(), phone, criteria)
			if err != nil {
				if firebase.IsUserNotFound(err) {
					utils.StdOutError("Not Found\t %s \t %s", phone, err.Error())
					continue
				}
				utils.StdOutError("Error\t%s\t%s", phone, err.Error())
				continue
			}
			// @todo expand this list of users
			utils.StdOutSuccess("%s\tWas successfully Retrieved\n", user.UID)
		}
		os.Exit(0)
	},
}

func init() {
	findUserCmd.AddCommand(byPhoneCmd)
}
