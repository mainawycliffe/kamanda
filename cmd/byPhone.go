package cmd

import (
	"context"
	"fmt"
	"os"

	fAuth "firebase.google.com/go/auth"
	"github.com/mainawycliffe/kamanda/firebase"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/mainawycliffe/kamanda/views"
	"github.com/spf13/cobra"
)

// byPhoneCmd represents the byPhone command
var byPhoneCmd = &cobra.Command{
	Use:     "byPhone",
	Aliases: []string{"phone", "by-phone"},
	Short:   "find a Firebase Auth User by their phone number",
	Example: `kamanda auth find by-phone +254712345678`,
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
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one Firebase user UID is required!")
			os.Exit(1)
		}
		// args = list of phone numbers
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one Firebase user UID is required!")
			os.Exit(1)
		}
		criteria := auth.ByUserUIDCriteria
		users := make([]*fAuth.ExportedUserRecord, 0)
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
		// draw table
		views.ViewUsersTable(users, "")
		os.Exit(0)
	},
}

func init() {
	findUserCmd.AddCommand(byPhoneCmd)
}
