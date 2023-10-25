package cmd

import (
	"fmt"
	"github.com/mstxq17/MoreFind/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the semantic version number of MoreFind",
	Long:  `All software has versions. This is MoreFind's`,
	Run: func(cmd *cobra.Command, args []string) {
		v := fmt.Sprintf("MoreFind %s", common.VERSION)
		fmt.Println(v)
	},
}
