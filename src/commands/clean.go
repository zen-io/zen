package commands

import (
	"github.com/baulos-io/baulos/src/engine"

	"github.com/spf13/cobra"
)

var cleanCmd = &AhoyCommand{
	Pre:  nil,
	Post: nil,
	Command: func(eng *engine.Engine) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "clean",
			Short:             "Clean package caches",
			Long:              `Clean the cache for the provided packages`,
			ValidArgsFunction: eng.AutocompleteTargets,
			Args:              cobra.MinimumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				eng.CleanCache(args)
			},
		}

		return cmd
	},
}
