package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/baulos-io/baulos/src/commands"
	"github.com/baulos-io/baulos/src/engine"
)

var eng *engine.Engine

func init() {
	var err error
	eng, err = engine.NewEngine()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rootCmd.PersistentFlags().IntP("verbosity", "v", 3, "Verbosity level")
	rootCmd.PersistentFlags().Bool("clean", false, "Run a with a clean config")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Execute a dry run")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().Bool("shell", false, "Enter a shell instead of executing")
	rootCmd.PersistentFlags().Bool("raw-output", false, "Enable raw output, disabling tty")
	rootCmd.PersistentFlags().Bool("keep-output", false, "Keep the output of finished tasks")

	newTargetCmd.Flags().StringP("name", "n", "", "Name for the target")
	newTargetCmd.MarkFlagRequired("name")
	newTargetCmd.Flags().StringP("type", "t", "", "Target type")
	newTargetCmd.MarkFlagRequired("type")

	newPkgCmd.Flags().String("path", "", "Path for the target")
	newPkgCmd.MarkFlagRequired("path")
	newPkgCmd.Flags().String("project", "", "Project where to create the target")
	newPkgCmd.MarkFlagRequired("type")
	newDescTargetCmd.Flags().StringP("type", "t", "", "Target type")
	newDescTargetCmd.MarkFlagRequired("name")
	loginCmd.Flags().StringP("profile", "p", "login", "Profile to use for logging in")

	newDescCmd.AddCommand(newDescTargetCmd)
	newNewCmd.AddCommand(newTargetCmd)
	newNewCmd.AddCommand(newPkgCmd)

	rootCmd.AddCommand(newDescCmd)
	rootCmd.AddCommand(newNewCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(loginCmd)

	fns := map[string]*engine.RunFnMap{}
	for _, cmd := range commands.ExportedCommands {
		resolvedCmd := cmd.Command(eng)
		rootCmd.AddCommand(resolvedCmd)
		fns[resolvedCmd.Name()] = &engine.RunFnMap{
			Pre:  cmd.Pre,
			Post: cmd.Post,
		}
	}

	eng.RegisterCommandFunctions(fns)
}

var rootCmd = &cobra.Command{
	Use:          "baulos",
	Short:        "Ahoy is a build and deploy system",
	Long:         `A fast, iterative build and deploy system`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "debug" {
			cmd.Flags().Set("raw-output", "true")
		}
		eng.Initialize(cmd.Flags())
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		eng.Done()
		return nil
	},
}

func Execute() {
	rootCmd.Execute()
}
