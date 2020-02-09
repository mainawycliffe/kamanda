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
	Short:   "Find a user by uid. To find user by email or by phone use `find byEmail` or `find byPhone`",
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError("atleast one Firebase user UID is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserUIDCriteria
		for _, uid := range args {
			user, err := auth.GetUser(context.Background(), uid, criteria)
			if err != nil {
				if firebase.IsUserNotFound(err) {
					utils.StdOutError("Not Found\t %s \t %s", uid, err.Error())
					continue
				}
				utils.StdOutError("Error \t %s \t %s", uid, err.Error())
				continue
			}
			//@todo something with the output
			utils.StdOutSuccess("Success \t %s \t Was successfully Retrieved", user.UID)
		}
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(findUserCmd)
	findUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
