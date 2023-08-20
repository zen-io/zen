package commands

import (
	"github.com/zen-io/zen-engine/engine"

	"github.com/spf13/cobra"
)

type ZenCommand struct {
	Pre             func(eng *engine.Engine, es *engine.ExecutionStep) error
	Post            func(eng *engine.Engine, es *engine.ExecutionStep) error
	Command         func(eng *engine.Engine) *cobra.Command
	SubCommands     []*ZenCommand
}

var ExportedCommands = []*ZenCommand{
	deployCmd, buildCmd, cleanCmd, removeCmd, runCmd, queryCmd,
}
