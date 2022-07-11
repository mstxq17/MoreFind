package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
	"mvdan.cc/xurls/v2"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func isIPAddr(domain string) bool {
	ipaddr := net.ParseIP(domain)
	return ipaddr != nil
}

func isPrivateIP(line string) bool {
	var iIRegex = regexp.MustCompile("^(10.\\d{1,3}.\\d{1,3}.((0/([89]|1[0-9]|2\\d|3[012]))|(\\d{1,3})))|(172.(1[6789]|2\\d|3[01]).\\d{1,3}.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)|(192.168.\\d{1,3}.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)$")
	return iIRegex.MatchString(line)
}

func searchUrl(line string) []string {
	rxStrict := xurls.Relaxed()
	return rxStrict.FindAllString(line, -1)
}

func searchDomain(line string, rootDomain bool) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "http") == false {
		line = "https://" + line
	}
	u, err := url.Parse(line)
	if err != nil {
		log.Fatal(err)
	}
	if rootDomain {
		return searchRootDomain(u.Hostname())
	} else {
		return u.Hostname()
	}
}

// Reference:https://pkg.go.dev/golang.org/x/net/publicsuffix
/*
Description: search the eTLD + 1(rootDOmain) from the completed domain
param domain: completed domain
return: rootDomain
*/
func searchRootDomain(domain string) string {
	eTLD, _ := publicsuffix.EffectiveTLDPlusOne(domain)
	return eTLD
}

func searchIp(line string) []string {
	// only support ipv4, ipv6 will be supported in future
	var ipRegex = regexp.MustCompile("((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))")
	return ipRegex.FindAllString(line, -1)
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
		log.Fatal("len Range Invalid, format should be 'min-max', ex 0-100")
		return 0, 0
	}
}

var (
	file         string
	output       string
	myUrl        bool
	myDomain     bool
	myRootDomain bool
	myIp         bool
	myPrivateIp  bool
	myLimitLen   string
	myShow       bool
	rootCmd      = &cobra.Command{
		Use:   "morefind",
		Short: "MoreFind is a very fast script for searching URL、Domain and Ip from specified stream",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
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
			r := bufio.NewReader(_file)
			// todo: current structure may be chaotic, should abstract the handle process
			// if show flag be selected，deal with it first
			if myUrl == false && myDomain == false && myIp == false {
				if myShow == true {
					count := 0
					maxLength := 0
					minLength := 0
					first := true
					for {
						line, err := r.ReadString('\n')
						if err != nil {
							break
						}
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
						fmt.Printf("%-5d%-7s\t%s", count, " Len:"+lineLength, line)
					}
					fmt.Println("\n==================================================")
					fmt.Printf("CountLine: %d MaxLength: %d, MinLength: %d\n", count, maxLength, minLength)
					return
				}
				if myLimitLen != "" {
					min, max := filterLen(myLimitLen)
					for {
						line, err := r.ReadString('\n')
						line = strings.TrimSpace(line)
						if err != nil {
							break
						}
						if min <= len(line) && len(line) <= max {
							fmt.Println(line)
						}
					}
					return
				}
			}
			if myUrl == false && myDomain == false && myIp == false {
				myUrl = true
			}
			var urlList []string
			var domainList []string
			var ipList []string
			// remove duplicate url
			found := make(map[string]bool)
			for {
				line, err := r.ReadString('\n')
				if err != nil && err != io.EOF {
					panic(err)
				}
				if err == io.EOF {
					break
				}
				line = strings.TrimSpace(line)
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
								fmt.Println(_url)
								found[_url] = true
							}
						}
						if myDomain == true {
							_domain := searchDomain(_url, myRootDomain)
							if _domain == "" || isIPAddr(_domain) {
								continue
							}
							if output != "" {
								domainList = append(domainList, _domain)
							}
							// remove repeated string
							if _, ok := found[_domain]; !ok {
								fmt.Println(_domain)
								found[_domain] = true
							}
						}
					}
				}
				if myIp == true {
					searchIp := searchIp(line)
					for _, _ip := range searchIp {
						if output != "" {
							ipList = append(ipList, _ip)
						}
						// remove repeated string
						if _, ok := found[_ip]; !ok {
							if myPrivateIp == true {
								if isPrivateIP(_ip) == false {
									fmt.Println(_ip)
									found[_ip] = true
								}
							} else {
								fmt.Println(_ip)
								found[_ip] = true
							}
						}
					}
				}
			}
			if output != "" {
				_output, err := os.Create(output)
				if err != nil {
					log.Fatal(err)
				}
				defer func(_output *os.File) {
					err := _output.Close()
					if err != nil {
						log.Fatal(err)
					}
				}(_output)
				writer := bufio.NewWriter(_output)
				for key := range found {
					_, err := writer.WriteString(key + "\n")
					if err != nil {
						return
					}
				}
				err = writer.Flush()
				if err != nil {
					log.Fatal(err)
					return
				}
			}
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
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "search the info in specified file")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output the result to specified file")
	rootCmd.PersistentFlags().BoolVarP(&myIp, "ip", "i", false, "search ip from stdin or file")
	rootCmd.PersistentFlags().BoolVarP(&myPrivateIp, "exclude", "", false, "exclude internal/private segment of ip when searching ip")
	rootCmd.PersistentFlags().BoolVarP(&myDomain, "domain", "d", false, "search domain from stdin or file")
	rootCmd.PersistentFlags().BoolVarP(&myRootDomain, "root", "", false, "only output the rootDomain when searching domain")
	rootCmd.PersistentFlags().BoolVarP(&myUrl, "url", "u", false, "search url from stdin or file")
	rootCmd.PersistentFlags().StringVarP(&myLimitLen, "len", "l", "", "search specify the length of string, \"-l 35\" == \"-l 0-35\" ")
	rootCmd.PersistentFlags().BoolVarP(&myShow, "show", "s", false, "show the length of each line and summaries")
}
