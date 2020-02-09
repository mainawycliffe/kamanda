package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"list", "listUsers"},
	Short:   "Get a list of users in firebase auth.",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("nextPageToken")
		if err != nil {
			utils.StdOutError("An error occurred while parsing next page token")
			os.Exit(1)
		}
		getUsers, err := auth.ListUsers(context.Background(), 0, token)
		if err != nil {
			utils.StdOutError("Error! %s", err.Error())
			os.Exit(1)
		}
		// @todo: do something with the response i.e save to file or list them
		utils.StdOutSuccess("Next Page Token %s", getUsers.NextPageToken)
		utils.StdOutSuccess("Number of Users: %d", len(getUsers.Users))
	},
}

func init() {
	authCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().StringP("nextPageToken", "n", "", "Fetch next set of results")
}
