package cmd

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/mstxq17/MoreFind/core"
	"github.com/mstxq17/MoreFind/update"
	"github.com/mstxq17/MoreFind/vars"
	"github.com/spf13/cobra"
	"golang.org/x/net/publicsuffix"
	"log"
	"mvdan.cc/xurls/v2"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	MaxTokenSize = 512 * 1024 * 1024
)

var logger *log.Logger

// IPAndPort define custom struct
// 自定义一个结构体
type IPAndPort struct {
	IP   string
	Port string
}

type ErrorCallback func() *log.Logger

func isIPAddr(domain string) bool {
	ipaddr := net.ParseIP(domain)
	return ipaddr != nil
}

func isPrivateIP(line string) bool {
	// update regex pattern to match loopback and private ip
	// 更新正则表达式模式以匹配环回和私有IP
	//var iIRegex = regexp.MustCompile("^(10.\\d{1,3}.\\d{1,3}.((0/([89]|1[0-9]|2\\d|3[012]))|(\\d{1,3})))|(172.(1[6789]|2\\d|3[01]).\\d{1,3}.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)|(192.168.\\d{1,3}.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)$")
	var iIRegex = regexp.MustCompile("^(127\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}|10\\.\\d{1,3}\\.\\d{1,3}\\.((0/([89]|1[0-9]|2\\d|3[012]))|(\\d{1,3}))|172\\.(1[6789]|2\\d|3[01])\\.\\d{1,3}\\.\\d{1,3}(/(1[6789]|2\\d|3[012]))?|192\\.168\\.\\d{1,3}\\.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)$")
	return iIRegex.MatchString(line)
}

func searchUrl(line string) []string {
	rxRelaxed := xurls.Relaxed()
	result := rxRelaxed.FindAllString(line, -1)
	return result
}

func searchDomain(line string, rootDomain bool) (string, string) {
	/**
	匹配域名并输出
	match domain format and output
	*/
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "http") == false {
		line = "https://" + line
	}
	// bug fix #2
	// 修复issue #2
	_, exists := os.LookupEnv("hiddenDev")
	if !exists {
		line = deepMakeUrl(line)
	}
	u, err := url.Parse(line)
	if err != nil {
		// 直接抛出错误
		logger.Println(err)
		return "", ""
	}
	domain := u.Hostname()
	port := u.Port()
	// match the domain strictly
	// 严格匹配域名格式
	index := strings.Index(domain, ",")
	// 修复存在逗号的bug
	// patch the bug it contains  comma
	if index >= 0 {
		domain = domain[:index]
	}
	if isIPAddr(domain) {
		return port, domain
	}
	if rootDomain {
		return port, searchRootDomain(domain)
	} else {
		return port, domain
	}
}

// Reference:https://pkg.go.dev/golang.org/x/net/publicsuffix
/*
Description: search the eTLD + 1(rootDomain) from the completed domain
param domain: completed domain
return: rootDomain
*/
func searchRootDomain(domain string) string {
	eTLD, _ := publicsuffix.EffectiveTLDPlusOne(domain)
	return eTLD
}

func searchIp(line string) []IPAndPort {
	// only support ipv4, ipv6 will be supported in future
	//var ipRegex = regexp.MustCompile("((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))")
	ipPortRegex := regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))(:\d{1,5})?`)
	matches := ipPortRegex.FindAllStringSubmatch(line, -1)
	// store entries of result
	// 保存多个结果
	var result []IPAndPort
	for _, match := range matches {
		ip := match[1]
		port := match[8]
		if port != "" {
			port = port[1:]
		}
		entry := IPAndPort{IP: ip, Port: port}
		result = append(result, entry)
	}
	return result
}

func filterLen(lenRange string) (int, int) {
	standardPattern := regexp.MustCompile("^\\d+-\\d+$")
	oneIntPattern := regexp.MustCompile("^\\d+$")
	if standardPattern.MatchString(lenRange) {
		splitRes := strings.Split(lenRange, "-")
		minLength, _ := strconv.Atoi(splitRes[0])
		maxLength, _ := strconv.Atoi(splitRes[1])
		return minLength, maxLength
	} else if oneIntPattern.MatchString(lenRange) {
		maxLength, _ := strconv.Atoi(lenRange)
		return 0, maxLength
	} else {
		logger.Fatal("len Range Invalid, format should be 'min-max', ex 0-100")
		return 0, 0
	}
}

// the below two function can be merged and optimized
// 下面两个函数可以根据运行结构，将url解析那一部分抽象出来统一调用
func filterExt(_url string, filterExts string) bool {
	fileExt := fileExt(_url)
	_exts := strings.Split(filterExts, ",")
	// for improve the filtering speed, reducing the comparative work，use map
	// 为了提高速度，减少比较，使用map来判断
	extMap := map[string]int{}
	for _, suffix := range _exts {
		// convert to lowercase uniformly
		// 统一小写
		suffix = strings.TrimSpace(suffix)
		suffix = strings.ToLower(suffix)
		extMap[suffix] = 1
	}
	if _, ok := extMap[fileExt]; ok {
		return true
	} else {
		return false
	}
}

func fileExt(_url string) string {
	// bug fix #2
	// 修复issue #2
	_, exists := os.LookupEnv("hiddenDev")
	if !exists {
		_url = deepMakeUrl(_url)
	}
	u, err := url.Parse(_url)
	if err != nil {
		// ignore the exception for preventing from blocking next line
		// 忽略异常防止阻塞下一行的处理
		//logger.Fatal(err)
	}
	part := strings.Split(u.Path, "/")
	fileName := part[len(part)-1]
	if strings.Contains(fileName, ".") {
		filePart := strings.Split(fileName, ".")
		// convert to lowercase
		// 统一转换为小写
		return strings.ToLower(filePart[len(filePart)-1])
	} else {
		return ""
	}
}

func handleStdin(file string) (*os.File, os.FileInfo) {
	var _file *os.File
	if file != "" {
		var err error
		_file, err = os.Open(file)
		if err != nil {
			panic(err)
		}
	} else {
		_file = os.Stdin
	}
	// use features to solve  whether  has input
	// 利用特性解决程序是否有输入的问题
	fi, _ := _file.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		logger.Println("No input found, exit ...")
		// optimize exit logic
		// 优化退出逻辑
		os.Exit(0)
	}
	return _file, fi
}

func updateCommand(cmd *cobra.Command, args []string) {
	callBackError := func() *log.Logger {
		return logger
	}
	if myUpdate {
		update.GetUpdateToolCallback(vars.TOOLNAME, vars.VERSION, callBackError)()
	}
}

func preCommand(cmd *cobra.Command, args []string) bool {
	// 输出
	outputchan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go syncOutput(&wg, outputchan)
	// if cidr flag be selected，deal with it first
	// 如果选择 cidr 参数，首先处理它
	if myCidr != "" && myCidr != "__pipe__" {
		err := core.GenIP(myCidr, outputchan)
		if err != nil {
			logger.Println(err)
		}
		close(outputchan)
		wg.Wait()
		return true
	} else {
		close(outputchan)
		wg.Wait()
		return false
	}
}

func runCommand(cmd *cobra.Command, args []string) {
	// unified data stream
	// 统一数据流
	_file, fi := handleStdin(file)
	// prevent memory leaking
	// 防止内存泄漏
	defer func() {
		if err := _file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	// define global reader of input
	// 定义全局输入读取流
	var scanner *bufio.Scanner
	if myProgress {
		bar := pb.Full.Start64(fi.Size())
		defer func() {
			bar.Finish()
		}()
		reader := bar.NewProxyReader(bufio.NewReader(_file))
		scanner = bufio.NewScanner(reader)
	} else {
		reader := bufio.NewReader(_file)
		scanner = bufio.NewScanner(reader)
	}
	buf := make([]byte, 0, 64*1024)
	// support maximum  512MB buffer every line & support set maximum size through env
	// 支持最大读取单行 512MB 大小 & 支持环境变量设置更大值
	scanner.Buffer(buf, core.GetEnvOrDefault("MaxTokenSize", MaxTokenSize))
	// 输出
	outputchan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go syncOutput(&wg, outputchan)
	// todo: current structure may be chaotic, should abstract the handle process
	if myCidr == "__pipe__" {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			err := core.GenIP(line, outputchan)
			if err != nil {
				logger.Println(err)
			}
		}
	}
	if myUrl == false && myDomain == false && myIp == false {
		if myShow == true {
			count := 0
			maxLength := 0
			minLength := 0
			first := true
			for scanner.Scan() {
				line := scanner.Text()
				lineLength := strconv.Itoa(len(line))
				if len(line) > maxLength {
					maxLength = len(line)
				}
				if len(line) > 0 && first == true {
					minLength = len(line)
					first = false
				}
				if len(line) < minLength && first == false {
					minLength = len(line)
				}
				count++
				outputLine := fmt.Sprintf("%-5d Len:%-6s\t%s", count, lineLength, line)
				outputchan <- outputLine
			}
			splitPadding := "==================================================="
			outputchan <- splitPadding
			summaryTotal := fmt.Sprintf("CountLine: %d MaxLength: %d, MinLength: %d", count, maxLength, minLength)
			outputchan <- summaryTotal
		}
		if myLimitLen != "" {
			min, max := filterLen(myLimitLen)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if min <= len(line) && len(line) <= max {
					outputchan <- line
				}
			}
		}
	}
	if myUrl == false && myDomain == false && myIp == false {
		myUrl = true
	}
	var urlList []string
	var domainList []string
	var ipList []string
	// remove duplicated url
	// 去除重复的url
	found := make(map[string]struct{})
	// define stream myself
	// 定义自己的输出流
	var outputBuffer *core.MyBuffer
	var customStringHandler core.CustomStringHandler
	if myRule != "" {
		outputBuffer = core.NewMyBuffer(true)
		customStringHandler.Strategy = 1
		customStringHandler.Rule = myRule
		customStringHandler.Flag = myFlag
	} else {
		outputBuffer = core.NewMyBuffer(false)
		customStringHandler.Strategy = 0
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if myUrl == true || myDomain == true {
			searchUrl := searchUrl(line)
			for _, _url := range searchUrl {
				_url = strings.TrimSpace(_url)
				if myUrl == true {
					if output != "" {
						urlList = append(urlList, _url)
					}
					// remove repeated string
					if _, ok := found[_url]; !ok {
						if myUrlFilter != "" {
							if !filterExt(_url, myUrlFilter) {
								outputBuffer.WriteString(_url, &customStringHandler)
								found[_url] = struct{}{}
							}
						} else {
							outputBuffer.WriteString(_url, &customStringHandler)
							found[_url] = struct{}{}
						}
						outputchan <- outputBuffer.TempString
					}
				}
				if myDomain == true {
					port, _domain := searchDomain(_url, myRootDomain)
					if _domain == "" || isIPAddr(_domain) {
						continue
					}
					if myWithPort {
						if port != "" {
							_domain = _domain + ":" + port
						}
					}
					if output != "" {
						domainList = append(domainList, _domain)
					}
					// remove repeated string
					if _, ok := found[_domain]; !ok {
						outputBuffer.WriteString(_domain, &customStringHandler)
						found[_domain] = struct{}{}
						outputchan <- outputBuffer.TempString
					}
				}
			}
		}
		if myIp == true {
			searchIp := searchIp(line)
			for _, ipps := range searchIp {
				ipWithPort := ipps.IP
				if myWithPort && ipps.Port != "" {
					ipWithPort = ipps.IP + ":" + ipps.Port
				}
				if output != "" {
					ipList = append(ipList, ipWithPort)
				}
				// remove repeated string
				// 删除重复的行
				if _, ok := found[ipWithPort]; !ok {
					if myPrivateIp == true {
						if isPrivateIP(ipWithPort) == false {
							outputBuffer.WriteString(ipWithPort, &customStringHandler)
							found[ipWithPort] = struct{}{}
						}
					} else {
						outputBuffer.WriteString(ipWithPort, &customStringHandler)
						found[ipWithPort] = struct{}{}
					}
					outputchan <- outputBuffer.TempString
				}
			}
		}
		outputBuffer.Reset()
	}
	// maybe exceed maxTokenSize length
	if err := scanner.Err(); err != nil {
		logger.Println(err)
	}
	close(outputchan)
	wg.Wait()
}

var (
	file         string
	output       string
	myUrl        bool
	myDomain     bool
	myRootDomain bool
	myWithPort   bool
	myIp         bool
	myPrivateIp  bool
	myLimitLen   string
	myShow       bool
	myUrlFilter  string
	myCidr       string
	myRule       string
	myFlag       string
	myProgress   bool
	myUpdate     bool
	myIPFormats  []string
	rootCmd      = &cobra.Command{
		Use:   "morefind",
		Short: "MoreFind is a very rapid script for extracting URL、Domain and Ip from data stream",
		Long:  "",

		Run: func(cmd *cobra.Command, args []string) {
			// run high priority command first
			// 先执行优先级高的命令,如额更新执行
			updateCommand(cmd, args)
			// 若 preCommand 返回 true，表示命令执行成功，直接返回
			if preCommand(cmd, args) {
				return
			}
			// 如果 preCommand 返回 false，继续执行 runCommand
			runCommand(cmd, args)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func init() {
	// set flag for global logger in init func
	// 在 init 函数中创建全局 logger 并设置标志
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	// reduce the  amount of calling function
	// 减少函数调用次数
	NewLine = core.NewLine()
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", vars.FileHelpEn)
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", vars.OutputHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myIp, "ip", "i", false, vars.IPHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myPrivateIp, "exclude", "", false, vars.ExcludeHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myDomain, "domain", "d", false, vars.DomainHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myRootDomain, "root", "", false, vars.RootDomainHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myWithPort, "port", "p", false, vars.WithPortHelpEn)
	rootCmd.PersistentFlags().StringVarP(&myRule, "rule", "r", "", vars.RuleHelpEn)
	rootCmd.PersistentFlags().StringVarP(&myFlag, "flag", "", "{}", vars.FlagHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myUrl, "url", "u", false, vars.URLHelpEn)
	rootCmd.PersistentFlags().StringVarP(&myUrlFilter, "filter", "", "", vars.URLFilterHelpEn)
	// this trick occurs from https://stackoverflow.com/questions/70182858/how-to-create-flag-with-or-without-argument-in-golang-using-cobra
	// help me a lot, so log it in the code， google dork: "flag needs an argument: cobra"
	// 感谢 https://stackoverflow.com/questions/70182858/how-to-create-flag-with-or-without-argument-in-golang-using-cobra 提供了如何解决--filter 默认参数的问题
	rootCmd.PersistentFlags().Lookup("filter").NoOptDefVal = "js,css,json,png,jpg,html,xml,zip,rar"
	rootCmd.PersistentFlags().StringVarP(&myCidr, "cidr", "c", "", vars.CidrHelpEn)
	rootCmd.PersistentFlags().StringSliceVarP(&myIPFormats, "alter", "a", nil, vars.AlterHelpEn)
	rootCmd.PersistentFlags().Lookup("cidr").NoOptDefVal = "__pipe__"
	rootCmd.PersistentFlags().StringVarP(&myLimitLen, "len", "l", "", vars.LimitLenHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myShow, "show", "s", false, vars.ShowHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myProgress, "metric", "m", false, vars.ProgressHelpEn)
	rootCmd.PersistentFlags().BoolVarP(&myUpdate, "update", "U", false, vars.UpdateHelpEn)
	// Dont sorted flag alphabetically
	// 禁止排序参数，按代码定义顺序展示
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false
}
