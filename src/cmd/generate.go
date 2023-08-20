package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zen-io/zen-engine/generator"
)

var newTargetCmd = &cobra.Command{
	Use:   "target",
	Short: "Generate a new target in a package",
	Long:  `Generate a target of name 'name' in the provided package`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// proj, _ := cmd.Flags().GetString("project")
		// t, _ := cmd.Flags().GetString("type")
		// generator.GenerateTarget(t, eng.PackageParser.KnownTypes(proj)[t])
		// return eng.BuildGraphAndRun(args, "pkg")
	},
}

var newPkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Generate a new package",
	Long:  `Generate a package in the provided path and project`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pkgPath, _ := cmd.Flags().GetString("path")
		project, _ := cmd.Flags().GetString("project")

		if eng.Projects[project] == nil {
			eng.Errorln("Project %s is not configured", project)
			return
		}
		path := filepath.Join(
			eng.Projects[project].Path,
			pkgPath,
			eng.Projects[project].Config.Parse.Filename,
		)

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			eng.Errorln("could not create path %s: %w", filepath.Dir(path), err)
			return
		}

		if file, err := os.Create(path); err != nil {
			eng.Errorln("could not create path %s: %w", filepath.Dir(path), err)
		} else {
			file.Close()
		}
	},
}

var newNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a template",
}

var newDescCmd = &cobra.Command{
	Use:   "desc",
	Short: "Describe a resource",
}

var newDescTargetCmd = &cobra.Command{
	Use:   "target-type",
	Short: "Describe a target type",
	Long:  `Describe a target type`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(generator.DescribeTarget(args[0]))
	},
}

// var newPkgCmd = &cobra.Command{
// 	Use:   "pkg",
// 	Short: "",
// 	Long:  `Build and deploy a list of packages and their dependencies`,
// 	Args:  cobra.MinimumNArgs(1),
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		return eng.BuildGraphAndRun(args, "pkg")
// 	},
// }

// var newPluginCmd = &cobra.Command{
// 	Use:   "plugin",
// 	Short: "",
// 	Long:  `Build and deploy a list of packages and their dependencies`,
// 	Args:  cobra.MinimumNArgs(1),
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		return eng.NewPluginRun(args, cmd.Flags())
// 	},
// }
