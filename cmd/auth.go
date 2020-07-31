package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:        "auth",
	Short:      "Perform Firebase Auth Operations",
	Deprecated: "it is being deprecated in favor or users command. Use \"kamanda users\" instead of \"kamanda auth\" users.",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
