package cmd

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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
		filePrepender := func(filename string) string {
			now := time.Now().Format(time.RFC3339)
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/commands/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}
		// Probably Rethink How This Works A Little Bit
		// Instead of Starting From Root, Probably Individual Categories is much
		// nice i.e. start with auth, then firestore, then storage etc
		err := doc.GenMarkdownTreeCustom(authCmd, "./docs/content/auth", filePrepender, linkHandler)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docsGenCmd)
}
