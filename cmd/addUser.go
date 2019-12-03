package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:     "addUser",
	Aliases: []string{"add-user", "add"},
	Short:   "Add a new Firebase user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addUser called")
	},
}

func init() {
	authCmd.AddCommand(addUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
