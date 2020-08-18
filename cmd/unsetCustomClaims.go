package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// unsetCustomClaimsCmd represents the unsetCustomClaims command
var unsetCustomClaimsCmd = &cobra.Command{
	Use:     "custom-claims",
	Aliases: []string{"customClaims"},
	Short:   "Unset custom claims from a firebase user account",
	Example: "kamanda users unset custom-claims --keys \"role\"",
	Run: func(cmd *cobra.Command, args []string) {
		keys, _ := cmd.Flags().GetStringArray("keys")
		hasError := false
		for _, v := range args {
			user, err := auth.RemoveCustomClaimFromUser(context.Background(), v, keys)
			if err != nil {
				hasError = true
				utils.StdOutError(os.Stderr, "Error removing custom claims for %s (%s)\n", v, user.Email)
				continue
			}
			utils.StdOutSuccess(os.Stdout, "Successfully remove custom claims for %s (%s)\n", v, user.Email)
		}
		if hasError {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	unsetCmd.AddCommand(unsetCustomClaimsCmd)
	unsetCustomClaimsCmd.Flags().StringArray("keys", nil, "Custom claims keys to remove")
	if err := unsetCustomClaimsCmd.MarkFlagRequired("keys"); err != nil {
		utils.StdOutError(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
