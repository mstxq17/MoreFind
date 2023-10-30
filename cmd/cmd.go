package cmd

import (
	"bufio"
	"fmt"
	"github.com/mstxq17/MoreFind/core"
	"github.com/mstxq17/MoreFind/update"
	"github.com/mstxq17/MoreFind/vars"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

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

var pattern string
var inverseMatch bool // Define a variable to hold the value of the inverse match flag

var grepCmd = &cobra.Command{
	Use:   "grep",
	Short: "If no grep , use this",
	Long:  `The grep command filters and displays lines matching a given pattern within files, akin to the Unix 'grep' command but without the second option.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			pattern = args[0]
		}
		fileStdin, _ := handleStdin(file)
		defer func() {
			if err := fileStdin.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		reader := bufio.NewReader(fileStdin)
		scanner := bufio.NewScanner(reader)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, MaxTokenSize)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			matchLine, err := core.MatchLine(line, pattern, inverseMatch)
			if err == nil && matchLine != "" {
				fmt.Println(matchLine)
			}
		}
	},
}

func init() {
	// try other style to parse params
	// 尝试使用不同的风格命令参数获取
	grepCmd.Flags().StringVarP(&pattern, "pattern", "P", "", "Pattern to search")
	grepCmd.Flags().BoolVarP(&inverseMatch, "inverse-match", "v", false, "Invert the match")
	grepCmd.SetUsageTemplate(usageTemplate)
	grepCmd.SetHelpTemplate(helpTemplate)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(grepCmd)
}
