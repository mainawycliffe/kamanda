package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/spf13/cobra"
)

// deleteUserCmd represents the deleteUser command
var deleteUserCmd = &cobra.Command{
	Use:   "deleteUser",
	Short: "Delete a Firebase User",
	Run: func(cmd *cobra.Command, args []string) {

		uid, err := cmd.Flags().GetString("uid")

		if err != nil {
			log.Fatalf("Error while reading uid flag: %s", err.Error())
		}

		fmt.Println(fmt.Sprintf("deleteUser called %v", uid))

		// call firebase auth to delete firebase
		err = auth.DeleteFirebaseUser(context.Background(), nil, uid)

		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	},
}

func init() {
	authCmd.AddCommand(deleteUserCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deleteUserCmd.Flags().String("uid", "", "The UID of the Firebase User you want to delete (required)")
	_ = deleteUserCmd.MarkFlagRequired("uid")
}
