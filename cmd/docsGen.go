package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsGenCmd represents the docsGen command
var docsGenCmd = &cobra.Command{
	Use:   "docsGen",
	Short: "Generate Kamanda Documentation",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docsGenCmd)
}
