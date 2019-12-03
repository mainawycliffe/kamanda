package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/mainawycliffe/kamanda/firebase"
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

		fmt.Println(fmt.Sprintf("deleteUser called %v", args))

		firebase := &firebase.Firebase{}

		err := firebase.InitializeFirbeaseApp(context.Background(), "")

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(deleteUserCmd)
}
