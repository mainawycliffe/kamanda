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
	Aliases: []string{"phone", "by-phone"},
	Short:   "find a Firebase Auth User by their phone number",
	Example: `kamanda auth find by-phone +254712345678`,
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of phone numbers
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one Firebase user UID is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserUIDCriteria
		for _, phone := range args {
			user, err := auth.GetUser(context.Background(), phone, criteria)
			if err != nil {
				if firebase.IsUserNotFound(err) {
					utils.StdOutError(os.Stderr, "Not Found\t %s \t %s", phone, err.Error())
					continue
				}
				utils.StdOutError(os.Stderr, "Error\t%s\t%s", phone, err.Error())
				continue
			}
			// @todo expand this list of users
			utils.StdOutSuccess(os.Stdout, "%s\tWas successfully Retrieved\n", user.UID)
		}
		os.Exit(0)
	},
}

func init() {
	findUserCmd.AddCommand(byPhoneCmd)
}
