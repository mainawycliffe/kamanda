package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd represents the users command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Perform Users Operations For Firebase User(s).",
}

func init() {
	authCmd.AddCommand(userCmd)
}
