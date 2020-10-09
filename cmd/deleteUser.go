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
	Use:     "delete",
	Aliases: []string{"delete-users", "deleteUsers"},
	Short:   "Delete multiple Firebase Auth User by their UID",
	Example: "kamanda users delete [uid1] [uid2] [uid3]",
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError(os.Stderr, "at least one Firebase user uid is required!")
			os.Exit(1)
		}
		// delete all listed user accounts
		for _, uid := range args {
			err := auth.DeleteFirebaseUser(context.Background(), uid)
			if err != nil {
				utils.StdOutError(os.Stderr, "%s - Not Deleted: %s\n", uid, err.Error())
				continue
			}
			utils.StdOutSuccess(os.Stdout, "%s - Deleted\n", uid)
		}
		os.Exit(0)
	},
}

func init() {
	usersCmd.AddCommand(deleteUserCmd)
}
