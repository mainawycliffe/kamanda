package cmd

import (
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set users details like password",
}

func init() {
	usersCmd.AddCommand(setCmd)
}
