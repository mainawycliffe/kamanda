package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// findUserCmd represents the findUser command
var findUserCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"findUser"},
	Short:   "Find a a Firebase Auth user by their Firebase UID.",
	Long: `Find a a Firebase Auth user by their Firebase UID.
To find user by email or by phone use "find by-email" or "find by-phone"`,
	Example: `kamanda auth find [UID1] [UID2]`,
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one Firebase user UID is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserUIDCriteria
		for _, uid := range args {
			user, err := auth.GetUser(context.Background(), uid, criteria)
			if err != nil {
				if firebase.IsUserNotFound(err) {
					utils.StdOutError(os.Stderr, "Not Found\t %s \t %s", uid, err.Error())
					continue
				}
				utils.StdOutError(os.Stderr, "Error \t %s \t %s", uid, err.Error())
				continue
			}
			//@todo something with the output
			utils.StdOutSuccess(os.Stdout, "Success \t %s \t Was successfully Retrieved", user.UID)
		}
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(findUserCmd)
	findUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
