package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
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
			utils.StdOutError("at least one Firebase user uid is required!")
			os.Exit(1)
		}
		// delete all listed user accounts
		for _, uid := range args {
			err := auth.DeleteFirebaseUser(context.Background(), uid)
			if err != nil {
				utils.StdOutError("%s - Not Deleted: %s\n", uid, err.Error())
				continue
			}
			utils.StdOutSuccess("%s - Deleted\n", uid)
		}
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(deleteUserCmd)
}
