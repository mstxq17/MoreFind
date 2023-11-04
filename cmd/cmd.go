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
var cmpMode int
var strictMode bool
var smart bool
var threshold int

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

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "File Comparator between two files",
	Long:  `File Comparator, a robust Golang tool, With options for strict or sorted comparison.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			file1Name := args[0]
			file2Name := args[1]
			lines1, err1 := core.ReadLines(file1Name)
			lines2, err2 := core.ReadLines(file2Name)
			if err1 != nil || err2 != nil {
				logger.Fatal(err1, err2)
			}
			onlyInA, onlyInB, inBoth := core.CompareFiles(lines1, lines2, strictMode)
			if cmpMode < 1 || cmpMode > 3 {
				logger.Fatalf("cmpMode value must between 1-3, you pass: %v", cmpMode)
			}
			if cmpMode == 1 {
				for _, line := range onlyInA {
					if line != "" {
						fmt.Println(line)
					}
				}
			}
			if cmpMode == 2 {
				for _, line := range onlyInB {
					if line != "" {
						fmt.Println(line)
					}
				}
			}
			if cmpMode == 3 {
				for _, line := range inBoth {
					if line != "" {
						fmt.Println(line)
					}
				}
			}
		} else {
			fmt.Println("Missing enough params ......")
			fmt.Printf("Usage: %v\t%s cmp a.txt b.txt -M [1/2/3]%v", core.NewLine(), vars.TOOLNAME, core.NewLine())
		}
	},
}

var deduCmd = &cobra.Command{
	Use:   "dedu",
	Short: "De-duplicated lines",
	Long:  `De-duplicated lines Applying multiple heuristics techniques`,
	Run: func(cmd *cobra.Command, args []string) {
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
		dr := core.NewDuplicateRemover(threshold, smart)
		for scanner.Scan() {
			line := scanner.Text()
			rResult := dr.RemoveDuplicator(line)
			if rResult != "" {
				fmt.Println(rResult)
			}
		}
	},
}

func init() {
	// try other style to parse params
	// 尝试使用不同的风格命令参数获取
	grepCmd.Flags().StringVarP(&pattern, "pattern", "P", "", vars.GrepPatternHelpEn)
	grepCmd.Flags().BoolVarP(&inverseMatch, "inverse-match", "v", false, vars.InverseMatchHelpEn)
	grepCmd.SetUsageTemplate(usageTemplate)
	grepCmd.SetHelpTemplate(helpTemplate)
	grepCmd.Flags().SortFlags = false
	// compare two file and match different mode result
	// 比较文件并匹配不同模式的结果
	diffCmd.Flags().IntVarP(&cmpMode, "mode", "M", 3, vars.DiffCmdHelpEn)
	diffCmd.Flags().BoolVarP(&strictMode, "strict", "", false, vars.StrictModeHelpEn)
	diffCmd.SetUsageTemplate(usageTemplate)
	diffCmd.SetHelpTemplate(helpTemplate)
	diffCmd.Flags().SortFlags = false
	// de-duplicated lines
	// 去重复行
	deduCmd.Flags().BoolVarP(&smart, "smart", "", false, vars.SmartHelpEn)
	deduCmd.Flags().IntVarP(&threshold, "threshold", "t", 15, vars.ThresholdHelpEn)
	deduCmd.SetUsageTemplate(usageTemplate)
	deduCmd.SetHelpTemplate(deduHelpTemplate)
	deduCmd.Flags().SortFlags = false
	// add to root command
	// 添加到 主命令
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(grepCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(deduCmd)
}
