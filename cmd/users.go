package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cheynewallace/tabby"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/mainawycliffe/kamanda/views"
	"github.com/spf13/cobra"
)

// usersCmd represents the listUsers command
var usersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"list", "listUsers", "list-users"},
	Short:   "Fetch and display a list of users in firebase auth.",
	Long: `This fetches users on Firebase Auth and either displays it in either simple table, an interactive console table, json or yaml format.
	
	By Default, this will display a simple table with the list of users.
	
	Use the output flag [--output] to display the data in either JSON or Yaml Formats. To show an interactive UI, please use the interactive flag [--interactive]. 
	
	This will display a table you can interact with.
	
In cases where there are more than 500 users, you will also get a nextPageToken, that you can use to fetch more users.`,
	Example: `kamanda users
kamanda users -o json
kamanda users -output yaml`,
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
		nextPageToken, err := cmd.Flags().GetString("nextPageToken")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading nextPageToken: %s", err.Error())
			os.Exit(1)
		}
		interactive, err := cmd.Flags().GetBool("interactive")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading interactive ui flag: %s", err.Error())
			os.Exit(1)
		}
		getUsers, err := auth.ListUsers(context.Background(), 0, nextPageToken)
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
		if interactive {
			// draw table
			views.ViewUsersTable(getUsers.Users, getUsers.NextPageToken)
			os.Exit(0)
		}
		// show minimal ui here
		header := []interface{}{"UID", "Email", "Display Name", "Provider", "Last Login", "Created On"}
		rows := make([][]interface{}, len(getUsers.Users))
		for index, user := range getUsers.Users {
			dateCreated := utils.FormatTimestampToDate(user.UserMetadata.CreationTimestamp, "02/01/2006 15:04:05 MST")
			lastModified := utils.FormatTimestampToDate(user.UserMetadata.LastLogInTimestamp, "02/01/2006 15:04:05 MST")
			row := []interface{}{
				user.UID,
				user.Email,
				user.DisplayName,
				user.ProviderID,

				lastModified,
				dateCreated,
			}
			rows[index] = row
		}
		views.SimpleTableList(tabby.New(), header, rows...).Print()
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.Flags().StringP("nextPageToken", "n", "", "Fetch next set of results")
	usersCmd.PersistentFlags().BoolP("interactive", "i", false, "Show Interactive UI for Users")
}
