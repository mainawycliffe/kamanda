package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/mainawycliffe/kamanda/views"
	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"list", "listUsers"},
	Short:   "Get a list of users in firebase auth.",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading output: %s", err.Error())
			os.Exit(1)
		}
		if output != "json" && output != "yaml" && output != "" {
			utils.StdOutError(os.Stderr, "Unsupported output!")
			os.Exit(1)
		}
		token, err := cmd.Flags().GetString("nextPageToken")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading nextPageToken: %s", err.Error())
			os.Exit(1)
		}
		getUsers, err := auth.ListUsers(context.Background(), 0, token)
		if err != nil {
			utils.StdOutError(os.Stderr, "Error! %s", err.Error())
			os.Exit(1)
		}
		formatedUsers, err := utils.FormatResults(getUsers, output)
		if err != nil && err.Error() != "Unknown Format" {
			utils.StdOutError(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		if formatedUsers != nil {
			fmt.Printf("%s\n", formatedUsers)
			os.Exit(0)
		}
		// draw table
		views.ViewUsersTable(getUsers.Users, getUsers.NextPageToken)
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().StringP("nextPageToken", "n", "", "Fetch next set of results")
}
