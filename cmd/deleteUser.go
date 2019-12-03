package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/spf13/cobra"
)

// deleteUserCmd represents the deleteUser command
var deleteUserCmd = &cobra.Command{
	Use:     "deleteUsers",
	Aliases: []string{"delete-users", "delete"},
	Short:   "Delete multiple Firebase Auth User by their UID",
	RunE: func(cmd *cobra.Command, args []string) error {
		// args = list of uids
		if len(args) == 0 {
			return errors.New("atleast one Firebase user uid is required!")
		}

		// delete all user accounts
		for _, uid := range args {

			err := auth.DeleteFirebaseUser(context.Background(), uid)

			if err != nil {
				return fmt.Errorf("An error occurred while deleting account with uid: %s :%w", uid, err)
			}
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(deleteUserCmd)
}
