package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// addUsersCmd represents the addUsers command
var addUsersCmd = &cobra.Command{
	Use:     "addUsers",
	Aliases: []string{"add-users"},
	Short:   "Add multiple users from file (JSON/YAML)",
	Run: func(cmd *cobra.Command, args []string) {
		sourceFile, _ := cmd.Flags().GetString("source")
		sourceFileExtension, _ := cmd.Flags().GetString("extension")
		file, err := os.Stat(sourceFile)
		if err != nil && !os.IsNotExist(err) {
			utils.StdOutError("Source file doesn't exist!")
			os.Exit(1)
		}
		if file.IsDir() {
			utils.StdOutError("%s Source file is a directory not folder!", sourceFile)
			os.Exit(1)
		}
		usersToCreate := make([]auth.NewUser, 0)
		err = utils.UnmashalFormatFile(sourceFile, sourceFileExtension, &usersToCreate)
		if err != nil {
			utils.StdOutError("Error decoding your config file: %s", err.Error())
			os.Exit(1)
		}
		failedAccountCreation := 0
		for _, v := range usersToCreate {
			userRecord, err := auth.NewFirebaseUser(context.Background(), &v)
			if err != nil {
				// @todo: unwrap the errors properly to better error messages
				utils.StdOutError("%s Failed - %s", v.Email, err.Error())
				failedAccountCreation++
				continue
			}
			utils.StdOutSuccess("✔✔ email: %s SUCCESS \t uid: %s \n", userRecord.Email, userRecord.UID)
			// @todo probably also add custom claims
		}
		if failedAccountCreation > 0 {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(addUsersCmd)
	addUsersCmd.Flags().StringP("source", "s", "", "file with list of users to create")
	if err := addUsersCmd.MarkFlagRequired("source"); err != nil {
		fmt.Print(aurora.Sprintf(aurora.Red("%s\n"), err.Error()))
		os.Exit(1)
	}
	addUsersCmd.Flags().StringP("extension", "e", "yaml", "Source file type - json or yaml")
	if err := addUsersCmd.MarkFlagRequired("extension"); err != nil {
		fmt.Print(aurora.Sprintf(aurora.Red("%s\n"), err.Error()))
		os.Exit(1)
	}
}
