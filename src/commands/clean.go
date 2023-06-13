package commands

import (
	"github.com/zen-io/zen-engine/engine"

	"github.com/spf13/cobra"
)

var cleanCmd = &ZenCommand{
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
