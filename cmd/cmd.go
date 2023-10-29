package cmd

import (
	"fmt"
	"github.com/mstxq17/MoreFind/update"
	"github.com/mstxq17/MoreFind/vars"
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
		v := fmt.Sprintf("MoreFind %s ", vars.VERSION)
		fmt.Print(v)
		latestVersion, err := update.GetLatestVersion(vars.TOOLNAME, vars.VERSION)
		if latestVersion != "" && err == nil {
			v := fmt.Sprintf(" -> But latest version %s has released, run with -U / --update get it", latestVersion)
			fmt.Println(v)
		}
		fmt.Println("")
	},
}
