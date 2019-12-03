package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addUsersCmd represents the addUsers command
var addUsersCmd = &cobra.Command{
	Use:     "addUsers",
	Aliases: []string{"add-users"},
	Short:   "Add multiple Firebase Auth Users from file - JSON or YAML Supported!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addUsers called")
	},
}

func init() {
	authCmd.AddCommand(addUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addUsersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addUsersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
