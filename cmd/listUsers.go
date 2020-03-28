package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mainawycliffe/kamanda/firebase/auth"
	"github.com/mainawycliffe/kamanda/utils"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"list", "listUsers"},
	Short:   "Get a list of users in firebase auth.",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("nextPageToken")
		if err != nil {
			utils.StdOutError(os.Stderr, "An error occurred while parsing next page token")
			os.Exit(1)
		}
		getUsers, err := auth.ListUsers(context.Background(), 0, token)
		if err != nil {
			utils.StdOutError(os.Stderr, "Error! %s", err.Error())
			os.Exit(1)
		}
		if getUsers.NextPageToken != "" {
			utils.StdOutSuccess(os.Stdout, "Next Page Token %s", getUsers.NextPageToken)
		}
		// @todo: set the width of the table to 100%
		app := tview.NewApplication()
		table := tview.NewTable()
		tableHeaderColumnNames := []string{"#", "UID", "Email", "Name", "Phone Number"}
		for index, value := range tableHeaderColumnNames {
			tableCell := tview.NewTableCell(value).SetTextColor(tcell.ColorWhite).SetBackgroundColor(tcell.ColorBlue)
			table.SetCell(0, index, tableCell)
		}
		for i, v := range getUsers.Users {
			isEmailVerified := ""
			if !v.EmailVerified {
				isEmailVerified = "(Not Verified)"
			}
			table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", i+1))).SetSelectable(false, false)
			table.SetCell(i+1, 1, tview.NewTableCell(v.UID)).SetSelectable(true, false)
			table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("%s (%s)", v.Email, isEmailVerified))).SetSelectable(true, false)
			table.SetCell(i+1, 3, tview.NewTableCell(v.DisplayName)).SetSelectable(true, false)
			table.SetCell(i+1, 4, tview.NewTableCell(v.PhoneNumber)).SetSelectable(true, false)
		}
		table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			// exit when the following keys ae pressed
			if key == tcell.KeyEscape || key == tcell.KeyCtrlW {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
			uid := table.GetCell(row, 1).Text
			fmt.Printf("%s\n", uid)
			// @todo: fetch from firebase the user details

		})
		// @todo add pagination to get more request
		if err := app.SetRoot(table, true).Run(); err != nil {
			utils.StdOutError(os.Stderr, "Error! %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	authCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().StringP("nextPageToken", "n", "", "Fetch next set of results")
}
