package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

var verifyEmailCmd = &cobra.Command{
	Use:     "isEmailVerified",
	Aliases: []string{"is-email-verified"},
	Short:   "Manually verify or un-verify a users email address.",
	Example: "kamanda users set isEmailVerified --status=false [UID] [UID]",
	Run: func(cmd *cobra.Command, args []string) {
		emailVerifiedStatus, _ := cmd.Flags().GetBool("status")
		hasError := false
		for _, v := range args {
			updatePassword := &auth.FirebaseUser{
				ShouldUpdateEmailVerified: true,
				EmailVerified:             emailVerifiedStatus,
			}
			user, err := auth.UpdateFirebaseUser(context.Background(), v, updatePassword)
			if err != nil {
				hasError = true
				utils.StdOutError(os.Stderr, "Error updating user %s email verified status\n", v)
				continue
			}
			utils.StdOutSuccess(os.Stdout, "Email verified status for user %s [%s] has been set to %s\n", user.UID, user.Email, emailVerifiedStatus)
		}
		if hasError {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	setCmd.AddCommand(verifyEmailCmd)
	verifyEmailCmd.Flags().Bool("status", false, "The email verified value to set to")
}
