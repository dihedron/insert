package cmd

import (
	"github.com/dihedron/put/text"
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
	replaceCmd.Flags().BoolP("once", "o", false, "Replace only the first occurrence")
}
