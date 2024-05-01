// Package cmd contains all the cobra commands for the CLI application.
package cmd

import (
	"errors"
	"strings"

	"github.com/89luca89/shell-funcheck/pkg/check"
	"github.com/89luca89/shell-funcheck/pkg/types"
	"github.com/spf13/cobra"
)

// NewCheckCommand will parse an input shell file and perform checks for undocumented
// global variables, arguments and outputs.
func NewCheckCommand() *cobra.Command {
	options := types.CLIOptions{}
	checkCommand := &cobra.Command{
		Use:   "check [options] [path-to-file]",
		Short: "Check for undocumented arguments and global variables inside functions",
		// RunE:             check,
		RunE: func(cmd *cobra.Command, _ []string) error {
			for _, filepath := range cmd.Flags().Args() {
				if !strings.HasPrefix(filepath, "--") {
					options.Files = append(options.Files, filepath)
				}
			}

			return runCheck(&options)
		},
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
	}

	checkCommand.Flags().SetInterspersed(false)

	checkCommand.Flags().BoolVar(&options.ExcludeComments, "exclude-comments", false, "")
	checkCommand.Flags().BoolVar(&options.WError, "werror", false, "")

	return checkCommand
}

func runCheck(options *types.CLIOptions) error {
	errored := false
	errlevel := 0

	// if we want werror, we increase sensitivity
	if options.WError {
		errlevel++
	}

	for _, filepath := range options.Files {
		vars, err := check.RunCheck(filepath, options)
		if err != nil {
			return err
		}

		for _, v := range vars {
			v.Print(filepath)

			if v.Level <= errlevel {
				errored = true
			}
		}
	}

	if errored {
		return errors.New("lint error")
	}

	return nil
}
