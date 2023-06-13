package commands

import (
	"github.com/baulos-io/baulos/src/engine"

	"github.com/spf13/cobra"
)

var deployCmd = &AhoyCommand{
	Pre:  nil,
	Post: nil,
	Command: func(eng *engine.Engine) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "deploy",
			Short:             "Deploy targets",
			Long:              `Build and deploy a list of packages and their dependencies`,
			Args:              cobra.MinimumNArgs(1),
			ValidArgsFunction: eng.AutocompleteTargets,
			Run: func(cmd *cobra.Command, args []string) {
				eng.CheckShellAndRun(cmd.Flags(), args, "deploy")
			},
		}

		cmd.Flags().StringP("env", "e", "", "Environment to manage")
		cmd.Flags().StringP("tag", "t", "", "Tag to manage")
		cmd.Flags().Bool("with-deps", false, "Also deploy dependencies of packages")
		cmd.Flags().Bool("parallel", false, "Ignore dependencies when deploying")
		cmd.Flags().Bool("clean", false, "Run a with a clean config")
		cmd.Flags().Bool("shell", false, "Run a shell debug session")

		return cmd
	},
}
