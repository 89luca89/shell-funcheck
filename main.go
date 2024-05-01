// Package main.
package main

import (
	"log"
	"strings"

	"github.com/89luca89/shell-funcheck/cmd"
	"github.com/89luca89/shell-funcheck/pkg/constants"
	"github.com/spf13/cobra"
)

func newApp() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "shell-funcheck",
		Short:            "Check shell functions for undocumented stuff",
		Version:          strings.TrimPrefix(constants.Version, "v"),
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
	}

	rootCmd.AddCommand(
		cmd.NewCheckCommand(),
	)

	return rootCmd
}

func main() {
	app := newApp()

	err := app.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
