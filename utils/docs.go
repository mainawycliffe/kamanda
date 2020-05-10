package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// DocsFrontMatter generates the needed front matter of Kamanda CLI options
func DocsFrontMatter(filename string, summary string) string {
	const fmTemplate = `---
title: "%s"
slug: %s
url: %s
summary: "%s"
---
`
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	url := "/commands/" + strings.ToLower(base) + "/"
	frontMatter := fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), base, url, summary)
	return frontMatter
}

// GenMarkdownTreeCustom is the same as the GenMarkdownTreeCustom from cobra.doc
// package but allows for use filePrepender function to accept a summary,
// usually the Short field of the command.
// Link: https://github.com/spf13/cobra/blob/master/doc/md_docs.md#customize-the-output
func GenMarkdownTreeCustom(cmd *cobra.Command, dir string, filePrepender func(string, string) string, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenMarkdownTreeCustom(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}
	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.WriteString(f, filePrepender(filename, cmd.Short)); err != nil {
		return err
	}
	if err := doc.GenMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}
