package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:    "auth",
	Short:  "Perform Firebase Auth Operations",
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
