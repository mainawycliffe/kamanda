package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/spf13/cobra"
)

// deleteUserCmd represents the deleteUser command
var deleteUserCmd = &cobra.Command{
	Use:     "deleteUsers",
	Aliases: []string{"delete-users", "delete"},
	Short:   "Delete multiple Firebase Auth User by their UID",
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			fmt.Printf("atleast one Firebase user uid is required!")
			os.Exit(1)
		}

		// delete all user accounts
		for _, uid := range args {

			err := auth.DeleteFirebaseUser(context.Background(), uid)

			if err != nil {
				fmt.Printf("Error deleting %s: %s\n", uid, err.Error())
				os.Exit(1)
			}
		}

		fmt.Printf("User was successfully deleted!")
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(deleteUserCmd)
}
