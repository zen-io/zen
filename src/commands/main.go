package commands

import (
	"github.com/baulos-io/baulos/src/cache"
	"github.com/baulos-io/baulosne"
	"github.com/baulos-io/bauloset"

	"github.com/spf13/cobra"
)

type AhoyCommand struct {
	Pre     func(eng *engine.Engine, target *target.Target, ci *cache.CacheItem) error
	Post    func(eng *engine.Engine, target *target.Target, ci *cache.CacheItem) error
	Command func(eng *engine.Engine) *cobra.Command
}

var ExportedCommands = []*AhoyCommand{
	deployCmd, buildCmd, cleanCmd, removeCmd, runCmd,
}
