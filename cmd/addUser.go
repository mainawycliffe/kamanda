package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/spf13/cobra"
)

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:     "addUser",
	Aliases: []string{"add-user", "add"},
	Short:   "Add a new Firebase user",
	Run: func(cmd *cobra.Command, args []string) {
		UID, _ := cmd.Flags().GetString("uid")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		emailVerified, _ := cmd.Flags().GetBool("emailVerified")
		phone, _ := cmd.Flags().GetString("phone")
		name, _ := cmd.Flags().GetString("name")
		photoURL, _ := cmd.Flags().GetString("photoURL")
		isDisabled, _ := cmd.Flags().GetBool("isDisabled")
		user := &auth.NewUser{
			UID:           UID,
			Email:         email,
			EmailVerified: emailVerified,
			PhoneNumber:   phone,
			DisplayName:   name,
			Password:      password,
			PhotoURL:      photoURL,
			Disabled:      isDisabled,
		}
		userRecord, err := auth.NewFirebaseUser(context.Background(), user)
		if err != nil {
			fmt.Print(aurora.Sprintf(aurora.Red("%s\n"), err.Error()))
			os.Exit(1)
		}
		// @todo: add custom claims in the future
		fmt.Print(aurora.Sprintf(aurora.Green("User created: %s email: %s\n"), userRecord.UID, userRecord.Email))
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().String("uid", "", "the uid of the user. autogenerated if absent")
	addUserCmd.Flags().String("email", "", "the email of the new user")
	addUserCmd.Flags().String("password", "", "the password of the new user")
	addUserCmd.Flags().Bool("emailVerified", false, "is the email verified")
	addUserCmd.Flags().String("phone", "", "the phone number of the user")
	addUserCmd.Flags().String("name", "", "the name of the user")
	addUserCmd.Flags().String("photoURL", "", "the photo url of the user")
	addUserCmd.Flags().Bool("isDisabled", false, "is the user account disabled")
	if err := addUserCmd.MarkFlagRequired("email"); err != nil {
		fmt.Print(aurora.Sprintf(aurora.Red("%s\n"), err.Error()))
		os.Exit(1)
	}
	if err := addUserCmd.MarkFlagRequired("password"); err != nil {
		fmt.Print(aurora.Sprintf(aurora.Red("%s\n"), err.Error()))
		os.Exit(1)
	}
}
