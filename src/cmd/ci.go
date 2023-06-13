package cmd

// import (
// 	"github.com/spf13/cobra"
// )

// var ciCmd = &cobra.Command{
// 	Use:   "ci",
// 	Short: "",
// 	Long:  `CI`,
// 	Args:  cobra.MinimumNArgs(1),
// }

// var ciTriggerCmd = &cobra.Command{
// 	Use:               "trigger",
// 	Short:             "Build packages",
// 	Long:              `Build a list of packages and their dependencies`,
// 	Args:              cobra.MinimumNArgs(1),
// 	ValidArgsFunction: AutocompleteTargets,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		eng.BuildGraphAndRun(args, "cover")
// 	},
// }

// var ciTailCmd = &cobra.Command{
// 	Use:               "tail",
// 	Short:             "Build packages",
// 	Long:              `Build a list of packages and their dependencies`,
// 	Args:              cobra.MinimumNArgs(1),
// 	ValidArgsFunction: AutocompleteTargets,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		eng.BuildGraphAndRun(args, "cover")
// 	},
// }

// var ciRunCmd = &cobra.Command{
// 	Use:               "tail",
// 	Short:             "Build packages",
// 	Long:              `Build a list of packages and their dependencies`,
// 	Args:              cobra.MinimumNArgs(1),
// 	ValidArgsFunction: AutocompleteTargets,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		eng.BuildGraphAndRun(args, "cover")
// 	},
// }
