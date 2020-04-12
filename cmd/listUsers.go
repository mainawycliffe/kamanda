package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/mainawycliffe/kamanda/views"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"list", "listUsers"},
	Short:   "Get a list of users in firebase auth.",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading output: %s", err.Error())
			os.Exit(1)
		}
		if output != "json" && output != "yaml" && output != "" {
			utils.StdOutError(os.Stderr, "Unsupported output!")
			os.Exit(1)
		}
		token, err := cmd.Flags().GetString("nextPageToken")
		if err != nil {
			utils.StdOutError(os.Stderr, "Error reading nextPageToken: %s", err.Error())
			os.Exit(1)
		}
		getUsers, err := auth.ListUsers(context.Background(), 0, token)
		if err != nil {
			utils.StdOutError(os.Stderr, "Error! %s", err.Error())
			os.Exit(1)
		}
		if output == "json" {
			json, err := json.Marshal(getUsers)
			if err != nil {
				utils.StdOutError(os.Stderr, "Error marsalling json: %s", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", json)
			os.Exit(0)
		}
		if output == "yaml" {
			yaml, err := yaml.Marshal(getUsers)
			if err != nil {
				utils.StdOutError(os.Stderr, "Error marsalling yaml: %s", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", yaml)
			os.Exit(0)
		}
		// draw table
		views.ViewUsersTable(getUsers.Users, getUsers.NextPageToken)
	},
}

func init() {
	authCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().StringP("nextPageToken", "n", "", "Fetch next set of results")
}
