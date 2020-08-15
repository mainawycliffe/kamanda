package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// enableUserAccountCmd represents the enableUserAccount command
var enableUserAccountCmd = &cobra.Command{
	Use:     "enabled",
	Short:   "Enable or disabled a user account",
	Example: "kamanda users set enabled --status=false [UID] [UID]",
	Run: func(cmd *cobra.Command, args []string) {
		isAccountEnabledAccount, _ := cmd.Flags().GetBool("status")
		hasError := false
		for _, v := range args {
			updatePassword := &auth.FirebaseUser{
				ShouldUpdateDisabled: true,
				Disabled:             !isAccountEnabledAccount,
			}
			user, err := auth.UpdateFirebaseUser(context.Background(), v, updatePassword)
			if err != nil {
				hasError = true
				if isAccountEnabledAccount {
					utils.StdOutError(os.Stderr, "Error occurred while enabling %s user account\n", v)
				}
				if !isAccountEnabledAccount {
					utils.StdOutError(os.Stderr, "Error occurred while disabling %s user account\n", v)
				}
				continue
			}
			if isAccountEnabledAccount {
				utils.StdOutSuccess(os.Stdout, "%s [%s] user account has been successfully enabled.\n", user.UID, user.Email)
			}
			if !isAccountEnabledAccount {
				utils.StdOutSuccess(os.Stdout, "%s [%s] user account has been successfully disabled.\n", user.UID, user.Email)
			}
		}
		if hasError {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	setCmd.AddCommand(enableUserAccountCmd)
	enableUserAccountCmd.Flags().Bool("status", false, "Status to toggle users accounts to.")
	if err := enableUserAccountCmd.MarkFlagRequired("status"); err != nil {
		utils.StdOutError(os.Stderr, "Status value to toggle user account to is required!\n")
		os.Exit(1)
	}
}
