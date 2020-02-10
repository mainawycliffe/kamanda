package cmd

import (
	"context"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

// customclaimsCmd represents the customclaims command
var customclaimsCmd = &cobra.Command{
	Use:     "customClaims",
	Aliases: []string{"claims", "custom-claims", "cc"},
	Short:   "Add custom claims to firebase user or users",
	Run: func(cmd *cobra.Command, args []string) {
		// args = list of uids
		if len(args) == 0 {
			utils.StdOutError("at least one Firebase user uid is required!")
			os.Exit(1)
		}
		customClaimsInput, _ := cmd.Flags().GetStringToString("customClaims")
		customClaims := utils.ProcessCustomClaimInput(customClaimsInput)
		for _, uid := range args {
			err := auth.AddCustomClaimToFirebaseUser(context.Background(), uid, customClaims)
			if err != nil {
				utils.StdOutError("%s - Failed: %s\n", uid, err.Error())
				continue
			}
			utils.StdOutSuccess("%s - Added Successfully\n", uid)
		}
		os.Exit(0)
	},
}

func init() {
	authCmd.AddCommand(customclaimsCmd)
	customclaimsCmd.Flags().StringToStringP("customClaims", "c", nil, "user custom claims i.e. --customClaims \"admin=true\"")
	if err := customclaimsCmd.MarkFlagRequired("customClaims"); err != nil {
		utils.StdOutError("%s\n", err.Error())
		os.Exit(1)
	}
}
