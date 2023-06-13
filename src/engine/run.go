package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/baulos-io/baulos/src/cache"

	"github.com/baulos-io/bauloset"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// special error that signals to stop the execution without errors
type DoNotContinue struct{}

type RunFnMap struct {
	Pre  func(eng *Engine, target *target.Target, ci *cache.CacheItem) error
	Post func(eng *Engine, target *target.Target, ci *cache.CacheItem) error
}

func (dnc DoNotContinue) Error() string {
	return "do not continue"
}

func (eng *Engine) _run_step(targetFqn string) error {
	fqn, err := target.NewFqnFromStr(targetFqn)
	if err != nil {
		return err
	}
	target := eng.targets[fqn.Qn()]
	script := fqn.Script()

	if target.TaskLogger, err = eng.out.CreateTask(targetFqn, ""); err != nil {
		return err
	}
	defer target.Done()

	// load cache
	var ci *cache.CacheItem
	ci, err = eng.Cache.LoadTargetCache(target)
	if err != nil {
		return fmt.Errorf("loading cache: %w", err)
	}

	if script == "build" {
		target.Cwd = ci.BuildCachePath()
	} else {
		target.Cwd = ci.BuildOutPath()
		if eng.Ctx.Env == "" {
			if len(target.Environments) == 1 {
				for e := range target.Environments {
					eng.Ctx.Env = e
					break
				}
			} else if len(target.Environments) > 0 {
				availableEnvs := []string{}
				for e := range target.Environments {
					availableEnvs = append(availableEnvs, e)
				}

				return fmt.Errorf("no environment was provided. Available options are %s", strings.Join(availableEnvs, ","))
			}
		}
		target.SetDeployVariables(
			eng.Ctx.Env,
			eng.ProjectDeployConfig(target.Project()).Variables,
			eng.config.Deploy.Variables,
		)
	}

	// pre run
	if eng.prePostFns[script] != nil && eng.prePostFns[script].Pre != nil {
		if err := eng.prePostFns[script].Pre(eng, target, ci); errors.Is(err, DoNotContinue{}) {
			return nil
		} else if err != nil {
			return fmt.Errorf("custom %s pre run: %w", script, err)
		}
	}

	if target.Scripts[script].Pre != nil {
		if err := target.Scripts[script].Pre(target, eng.Ctx); errors.Is(err, DoNotContinue{}) {
			return nil
		} else if err != nil {
			return fmt.Errorf("target %s pre run: %w", script, err)
		}
	}

	// run
	target.Env["CWD"] = target.Cwd
	if err := target.Scripts[script].Run(target, eng.Ctx); err != nil {
		target.Errorln("executing run: %s", err)
		return err
	}

	// POST RUN
	// custom script post run
	if eng.prePostFns[script] != nil && eng.prePostFns[script].Post != nil {
		if err := eng.prePostFns[script].Post(eng, target, ci); err != nil {
			return fmt.Errorf("custom %s post run: %w", script, err)
		}
	}

	// target post run
	if target.Scripts[script].Post != nil {
		if err := target.Scripts[script].Post(target, eng.Ctx); err != nil {
			return fmt.Errorf("target %s post run: %w", script, err)
		}
	}

	eng.Debugln("Finished %s", targetFqn)

	eng.out.CompleteTask(targetFqn)
	return nil
}

func (eng *Engine) CheckShellAndRun(flags *pflag.FlagSet, args []string, script string) {
	shell, _ := flags.GetBool("shell")
	if shell && len(args) > 1 {
		eng.Errorln("when using --shell, you can pass only one target")
		return
	}

	if err := eng.BuildGraph(args, script); err != nil {
		eng.Errorln("building graph: %w", err)
		return
	}

	if shell {
		fqn, err := target.NewFqnFromStr(args[0])
		if err != nil {
			eng.Errorln("inferring target fqn %s", args[0])
			return
		}
		fqn.SetDefaultScript(script)
		EnterTargetShell(eng.targets[fqn.Qn()], fqn.Script())
	}

	if err := eng.Run(); err != nil {
		eng.Errorln("executing the graph: %w", err)
	}
}

func (eng *Engine) AutocompleteTargets(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	targets, err := eng.AutocompleteTarget(toComplete)
	if err != nil {
		// fmt.Println(err)
		panic(err)
	}

	return targets, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
