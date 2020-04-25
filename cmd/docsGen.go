package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mainawycliffe/kamanda/utils"
	"github.com/spf13/cobra"
)

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

// docsGenCmd represents the docsGen command
var docsGenCmd = &cobra.Command{
	Use:   "docsGen",
	Short: "Generate Kamanda Documentation",
	Run: func(cmd *cobra.Command, args []string) {
		filePrepender := utils.DocsFrontMatter
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}
		// Probably Rethink How This Works A Little Bit
		// Instead of Starting From Root, Probably Individual Categories is much
		// nice i.e. start with auth, then firestore, then storage etc

		baseCommands := []struct {
			dir string
			cmd *cobra.Command
		}{
			{"", versionCmd},
			{"", loginCICmd},
			{"", loginCmd},
			{"", logoutCmd},
			{"/auth", authCmd},
		}

		for _, v := range baseCommands {
			dir := fmt.Sprintf("./docs/content/commands%s", v.dir)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				e := os.Mkdir(dir, os.ModePerm)
				if e != nil {
					utils.StdOutError(os.Stderr, "Error creating command directory: %s\n", e.Error())
					os.Exit(1)
				}
			}
			err := utils.GenMarkdownTreeCustom(v.cmd, dir, filePrepender, linkHandler)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(docsGenCmd)
}
