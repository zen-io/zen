package commands

import (
	"fmt"
	"os"

	"github.com/zen-io/zen-engine/engine"

	"github.com/spf13/cobra"
)

var deployCmd = &ZenCommand{
	Pre:  preDeploy,
	Post: postDeploy,
	Command: func(eng *engine.Engine) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "deploy",
			Short:             "Deploy targets",
			Long:              `Build and deploy a list of packages and their dependencies`,
			Args:              cobra.MinimumNArgs(1),
			ValidArgsFunction: eng.AutocompleteTargets,
			PreRun: func(cmd *cobra.Command, args []string) {
				eng.Ctx.UseEnvironments = true
			},
			Run: func(cmd *cobra.Command, args []string) {
				eng.ParseArgsAndRun(cmd.Flags(), args, "deploy")
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

var preDeploy = func(eng *engine.Engine, es *engine.ExecutionStep) error {
	if _, ok := es.Cache.Metadata["deploy"]; !ok {
		return nil
	}

	for _, v := range es.Cache.Metadata["deploy"].([]interface{}) {
		if v.(string) == eng.Ctx.Env {
			if err := es.Cache.ExpandOuts(es.Target.Outs); err != nil {
				if os.IsNotExist(err) {
					return nil
				}

				return fmt.Errorf("expanding outs: %w", err)
			}

			es.Target.SetStatus(fmt.Sprintf("Deployed %s [CACHED]", es.Target.Qn()))
			return engine.DoNotContinue{}
		}
	}
	return nil
}

var postDeploy = func(eng *engine.Engine, es *engine.ExecutionStep) error {
	if _, ok := es.Cache.Metadata["deploy"]; !ok {
		es.Cache.Metadata["deploy"] = []string{eng.Ctx.Env}
	} else {
		for _, v := range es.Cache.Metadata["deploy"].([]interface{}) {
			if v.(string) == eng.Ctx.Env {
				es.Cache.Metadata["deploy"] = append(es.Cache.Metadata["deploy"].([]string), eng.Ctx.Env)
			}
		}
	}

	if err := es.Cache.ExpandOuts(es.Target.Outs); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		
		return fmt.Errorf("expanding outs: %w", err)
	}

	es.Target.SetStatus(fmt.Sprintf("Deployed %s [DONE]", es.Target.Qn()))

	return nil
}
