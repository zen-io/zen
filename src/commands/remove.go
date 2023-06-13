package commands

import (
	"github.com/zen-io/zen-engine/engine"

	"github.com/spf13/cobra"
)

var removeCmd = &ZenCommand{
	Pre:  nil,
	Post: nil,
	Command: func(eng *engine.Engine) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "remove",
			Short:             "Remove targets",
			Long:              `Build and deploy a list of packages and their dependencies`,
			Args:              cobra.MinimumNArgs(1),
			ValidArgsFunction: eng.AutocompleteTargets,
			Run: func(cmd *cobra.Command, args []string) {
				eng.CheckShellAndRun(cmd.Flags(), args, "remove")
			},
		}

		cmd.Flags().StringP("env", "e", "", "Environment to manage")
		cmd.Flags().Bool("clean", false, "Run a with a clean config")

		return cmd
	},
}
