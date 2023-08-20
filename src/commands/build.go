package commands

import (
	"fmt"
	"os"

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
				eng.ParseArgsAndRun(cmd.Flags(), args, "build")
			},
		}

		return cmd
	},
}

var preBuild = func(eng *engine.Engine, es *engine.ExecutionStep) error {
	es.Target.SetStatus("Building %s", es.Target.Qn())

	// var ci *cache.CacheItem
	// ci, err := eng.Projects[es.Target.Project()].Cache.LoadTargetCache(es.Target, es.ExternalPath != nil, filepath.Dir(eng.Projects[es.Target.Project()].Config.PathForPackage(es.Target.Package())))
	// if err != nil {
	// 	return fmt.Errorf("loading cache: %w", err)
	// }
	// es.Cache = ci

	if !es.Clean && es.Cache.CheckCacheHits() {
		if err := es.Cache.ExpandOuts(es.Target.Outs); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("expanding outs: %w", err)
			}

			es.Target.Debugln("out not there: %w", err)
		} else {
			es.Cache.LoadMetadata()
			es.Target.SetStatus(fmt.Sprintf("Built %s [CACHED]", es.Target.Qn()))
			return engine.DoNotContinue{}
		}

		es.Target.Debugln("outs not present, building")
	} else {
		es.Target.Debugln("cache does not hit: %s", es.Cache.MetadataPath)
	}

	if err := es.Cache.DeleteCache(); err != nil {
		return fmt.Errorf("refreshing cache directory: %w", err)
	}

	if err := os.MkdirAll(es.Cache.BuildCachePath(), os.ModePerm); err != nil {
		return err
	}

	// Copy srcs to gen
	if err := es.Cache.LinkSrcsToCache(); err != nil {
		return fmt.Errorf("copying srcs into cache: %w", err)
	}

	// execute the custom build function
	return nil
}

var postBuild = func(eng *engine.Engine, es *engine.ExecutionStep) error {
	// copy outs to out dir

	if err := es.Cache.ExpandOuts(es.Target.Outs); err != nil {
		return fmt.Errorf("expanding outs: %w", err)
	}

	if err := es.Cache.LinkOutsIntoOut(); err != nil {
		return fmt.Errorf("copying outs after run: %w", err)
	}

	if es.Binary {
		for _, toBase := range es.Target.Outs {
			err := os.Chmod(toBase, 0755)
			if err != nil {
				return fmt.Errorf("making the target binary: %w", err)
			}
			break
		}
	}

	if slices.Contains(es.Target.Labels, "codegen") {
		if err := es.Cache.ExportOutsToPath(); err != nil {
			return fmt.Errorf("copying outs back to target path: %w", err)
		}
	}

	es.Target.SetStatus(fmt.Sprintf("Built %s [DONE]", es.Target.Qn()))

	return nil
}
