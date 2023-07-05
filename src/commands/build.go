package commands

import (
	"fmt"
	"os"

	"github.com/zen-io/zen-core/target"
	"github.com/zen-io/zen-engine/cache"
	"github.com/zen-io/zen-engine/engine"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var buildCmd = &ZenCommand{
	Pre:  preBuild,
	Post: postBuild,
	Command: func(eng *engine.Engine) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "build",
			Short:             "Build targets",
			Long:              `Build a list of packages and their dependencies`,
			Args:              cobra.MinimumNArgs(1),
			ValidArgsFunction: eng.AutocompleteTargets,
			Run: func(cmd *cobra.Command, args []string) {
				eng.CheckShellAndRun(cmd.Flags(), args, "build")
			},
		}

		return cmd
	},
}

var preBuild = func(eng *engine.Engine, target *target.Target, ci *cache.CacheItem) error {
	target.SetStatus("Building %s", target.Qn())

	if err := target.ExpandTools(eng.Cache.TargetOuts); err != nil {
		return fmt.Errorf("expanding tools: %w", err)
	}

	if err := target.InterpolateMyself(eng.Ctx); err != nil {
		return fmt.Errorf("interpolating myself: %w", err)
	}

	if !target.ShouldClean() && ci.CheckCacheHits() {
		if err := ci.ExpandOuts(target.Outs); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("expanding outs: %w", err)
			}
		} else {
			target.SetStatus(fmt.Sprintf("Built %s [CACHED]", target.Qn()))
			return engine.DoNotContinue{}
		}

		target.Debugln("outs not present, building")
	} else {
		target.Debugln("cache does not hit, building")
		if err := ci.DeleteCache(); err != nil {
			return fmt.Errorf("cleaning cache: %w", err)
		}
	}
	target.Cwd = ci.BuildCachePath()
	target.Env["PWD"] = target.Cwd

	if err := ci.DeleteCache(); err != nil {
		return fmt.Errorf("refreshing cache directory: %w", err)
	}

	if err := os.MkdirAll(ci.BuildCachePath(), os.ModePerm); err != nil {
		return err
	}

	// Copy srcs to tmp
	if err := ci.CopySrcsToCache(); err != nil {
		return fmt.Errorf("copying srcs into cache: %w", err)
	}

	// execute the custom build function
	return nil
}

var postBuild = func(eng *engine.Engine, target *target.Target, ci *cache.CacheItem) error {
	// copy outs to out dir

	if err := ci.ExpandOuts(target.Outs); err != nil {
		return fmt.Errorf("expanding outs: %w", err)
	}

	if err := ci.CopyOutsIntoOut(); err != nil {
		return fmt.Errorf("copying outs after run: %w", err)
	}

	if target.Binary {
		for _, toBase := range target.Outs {
			err := os.Chmod(toBase, 0755)
			if err != nil {
				return fmt.Errorf("making the target binary: %w", err)
			}
			break
		}
	}

	if slices.Contains(target.Labels, "codegen") {
		if err := ci.ExportOutsToPath(); err != nil {
			return fmt.Errorf("copying outs back to target path: %w", err)
		}
	}

	if err := ci.SaveMetadata(); err != nil {
		return fmt.Errorf("writing metadata: %w", err)
	}

	target.SetStatus(fmt.Sprintf("Built %s [DONE]", target.Qn()))

	return nil
}
