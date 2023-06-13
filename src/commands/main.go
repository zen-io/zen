package commands

import (
	"github.com/zen-io/zen-core/target"
	"github.com/zen-io/zen-engine/cache"
	"github.com/zen-io/zen-engine/engine"

	"github.com/spf13/cobra"
)

type ZenCommand struct {
	Pre         func(eng *engine.Engine, target *target.Target, ci *cache.CacheItem) error
	Post        func(eng *engine.Engine, target *target.Target, ci *cache.CacheItem) error
	Command     func(eng *engine.Engine) *cobra.Command
	SubCommands []*ZenCommand
}

var ExportedCommands = []*ZenCommand{
	deployCmd, buildCmd, cleanCmd, removeCmd, runCmd, queryCmd,
}
