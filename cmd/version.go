package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// VERSION number
// 版本号
const VERSION string = "1.4.4"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the semantic version number of MoreFind",
	Long:  `All software has versions. This is MoreFind's`,
	Run: func(cmd *cobra.Command, args []string) {
		v := fmt.Sprintf("MoreFind %s", VERSION)
		fmt.Println(v)
	},
}
