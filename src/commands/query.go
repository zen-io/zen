package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zen-io/zen-engine/engine"
)

var queryTargetsCmd = &ZenCommand{
	Command: func(eng *engine.Engine) *cobra.Command {
		return &cobra.Command{
			Use:               "targets",
			Short:             "Query targets",
			Long:              `Query all targets in a package tree`,
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: eng.AutocompleteTargets,
			Run: func(cmd *cobra.Command, args []string) {
				ts, err := eng.ExpandTargets(args, "")
				if err != nil {
					return
				}
				for _, t := range ts {
					fmt.Println(t[0 : len(t)-1])
					// eng.GetAllTargetsInPackage(t)
				}
			},
		}
	},
}

var queryCmd = &ZenCommand{
	Command: func(eng *engine.Engine) *cobra.Command {
		return &cobra.Command{
			Use:   "query",
			Short: "Query the packages",
		}
	},
	SubCommands: []*ZenCommand{queryTargetsCmd},
}
