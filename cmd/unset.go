package cmd

import (
	"github.com/spf13/cobra"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "A base command for unsetting user account data",
}

func init() {
	usersCmd.AddCommand(unsetCmd)
}
