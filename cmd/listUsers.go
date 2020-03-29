package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		tableHeaderColumnNames := []string{"#", "UID", "Email", "Email Verified", "Name", "Phone Number", "Provider"}
		for index, value := range tableHeaderColumnNames {
			tableCell := tview.NewTableCell(value).
				SetTextColor(tcell.ColorBlue).
				SetSelectable(false).
				SetStyle(tcell.StyleDefault).
				SetAlign(tview.AlignLeft)
			table.SetCell(0, index, tableCell)
		}
		for index, user := range getUsers.Users {
			isEmailVerified := "Yes"
			if !user.EmailVerified {
				isEmailVerified = "No"
			}
			table.SetCell(index+1, 0, tview.NewTableCell(fmt.Sprintf("%d", index+1))).SetSelectable(false, false)
			table.SetCell(index+1, 1, tview.NewTableCell(user.UID).SetExpansion(3)).SetSelectable(true, false)
			table.SetCell(index+1, 2, tview.NewTableCell(user.Email).SetExpansion(2)).SetSelectable(true, false)
			table.SetCell(index+1, 3, tview.NewTableCell(isEmailVerified).SetExpansion(1)).SetSelectable(true, false)
			table.SetCell(index+1, 4, tview.NewTableCell(user.DisplayName).SetExpansion(3)).SetSelectable(true, false)
			table.SetCell(index+1, 5, tview.NewTableCell(user.PhoneNumber).SetExpansion(1)).SetSelectable(true, false)
			table.SetCell(index+1, 6, tview.NewTableCell(strings.ToUpper(user.ProviderID)).SetExpansion(1)).SetSelectable(true, false)
		}
		table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			// exit when the following keys ae pressed
			if key == tcell.KeyEscape || key == tcell.KeyCtrlW {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).
			// @todo: show more info about the user when a row is clicked
			SetSelectedFunc(func(row int, column int) {
				// @todo: fetch from firebase the user details
				// uid := table.GetCell(row, 1).Text
			})
		tableContainer := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(table, 0, 1, true)
		appContainer := tview.NewFlex().AddItem(tableContainer, 0, 1, true)
		// @todo add pagination to get more request
		if err := app.SetRoot(appContainer, true).Run(); err != nil {
			utils.StdOutError(os.Stderr, "Error! %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	authCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().StringP("nextPageToken", "n", "", "Fetch next set of results")
}
