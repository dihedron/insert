package cmd

import (
	"github.com/dihedron/zed/text"
	"github.com/spf13/cobra"
)

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "provides line-by-line replacement",
	Long: `
	Replace substitutes the lines that match the gven pattern with
	the given, parametric text; capturing groups are supported via
	positional variables of the form ${<position>}, e.g. ${0}.`,
	Args: cobra.RangeArgs(3, 5),
	Run:  text.Replace,
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// replaceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// replaceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
