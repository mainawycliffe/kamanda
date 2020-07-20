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

// updateUserCmd represents the updateUser command
var updateUserCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"updateUser", "update-user"},
	Short:   "Update User Details",
	Long:    `Updates an existing Firebase users details given the UID`,
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "At least one user uid is required!")
			os.Exit(1)
		}
		if len(args) > 1 {
			utils.StdOutError(os.Stderr, "You can only update one user at a time")
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading output: %s", err.Error())
			os.Exit(1)
		}
		if output != "json" && output != "yaml" && output != "" {
			utils.StdOutError(os.Stderr, "Unsupported output!")
			os.Exit(1)
		}
		UID := args[0]
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		emailVerified, _ := cmd.Flags().GetBool("emailVerified")
		phone, _ := cmd.Flags().GetString("phone")
		name, _ := cmd.Flags().GetString("name")
		photoURL, _ := cmd.Flags().GetString("photoURL")
		isDisabled, _ := cmd.Flags().GetBool("isDisabled")
		user := &auth.FirebaseUser{
			UID:           UID,
			Email:         email,
			EmailVerified: emailVerified,
			PhoneNumber:   phone,
			DisplayName:   name,
			Password:      password,
			PhotoURL:      photoURL,
			Disabled:      isDisabled,
		}
		userRecord, err := auth.UpdateFirebaseUser(context.Background(), UID, user)
		if err != nil {
			utils.StdOutError(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		formatedUsers, err := utils.FormatResults(userRecord, output)
		if err != nil && err.Error() != "Unknown Format" {
			utils.StdOutError(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		if formatedUsers != nil {
			fmt.Printf("%s\n", formatedUsers)
			os.Exit(0)
		}
		// this is when no output is specified
		utils.StdOutSuccess(os.Stdout, "The following users have been updated successfully\n")
		header := []interface{}{"UID", "Email", "Email Verified", "Display Name", "Phone", "Disabled"}
		row := []interface{}{
			userRecord.UID,
			userRecord.Email,
			userRecord.EmailVerified,
			userRecord.DisplayName,
			userRecord.PhoneNumber,
			userRecord.Disabled,
		}
		views.SimpleTableList(tabby.New(), header, row).Print()
		os.Exit(0)
	},
}

func init() {
	listUsersCmd.AddCommand(updateUserCmd)
	updateUserCmd.Flags().String("email", "", "the email of the user (Required)")
	updateUserCmd.Flags().String("password", "", "the password of user (Required)")
	updateUserCmd.Flags().Bool("emailVerified", false, "is the email verified")
	updateUserCmd.Flags().String("phone", "", "the phone number of the user")
	updateUserCmd.Flags().String("name", "", "the name of the user")
	updateUserCmd.Flags().String("photoURL", "", "the photo url of the user")
	updateUserCmd.Flags().Bool("isDisabled", false, "is the user account disabled")
}
