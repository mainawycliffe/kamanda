package cmd

import (
	"context"
	"fmt"
	"os"

	fAuth "firebase.google.com/go/auth"
	"github.com/cheynewallace/tabby"
	"github.com/mainawycliffe/kamanda/firebase"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/mainawycliffe/kamanda/views"
	"github.com/spf13/cobra"
)

// byEmailCmd represents the byEmail command
var byEmailCmd = &cobra.Command{
	Use:     "byEmail",
	Aliases: []string{"email", "by-email"},
	Short:   "Find a Firebase Auth user by email address",
	Example: `kamanda auth find by-email email@example.com`,
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
		minimalUI, err := cmd.Flags().GetBool("minimalUI")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading minimal ui flag: %s", err.Error())
			os.Exit(1)
		}
		// args = list of UIDs
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one email is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserEmailCriteria
		users := make([]*fAuth.ExportedUserRecord, 0)
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
			users = append(users, &fAuth.ExportedUserRecord{
				UserRecord: user,
			})
		}
		formatedUsers, err := utils.FormatResults(users, output)
		if err != nil && err.Error() != "Unknown Format" {
			utils.StdOutError(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		if formatedUsers != nil {
			fmt.Printf("%s\n", formatedUsers)
			os.Exit(0)
		}
		if !minimalUI {
			// draw table
			views.ViewUsersTable(users, "")
			os.Exit(0)
		}
		// show minimal ui here
		header := []interface{}{"UID", "Email", "Display Name", "Provider", "Last Login", "Created On"}
		rows := make([][]interface{}, len(users))
		for index, user := range users {
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
	findUserCmd.AddCommand(byEmailCmd)
}
