package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mvdan.cc/xurls/v2"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func searchUrl(line string) []string {
	rxStrict := xurls.Relaxed()
	return rxStrict.FindAllString(line, -1)
}

func searchDomain(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "http") == false {
		line = "https://" + line
	}
	u, err := url.Parse(line)
	if err != nil {
		log.Fatal(err)
	}
	return u.Hostname()
}

func searchIp(line string) []string {
	// only support ipv4, ipv6 will be supported in future
	var ipRegex = regexp.MustCompile("((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))")
	return ipRegex.FindAllString(line, -1)
}

var (
	file     string
	output   string
	myUrl    bool
	myDomain bool
	myIp     bool
	rootCmd  = &cobra.Command{
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
			//sc := bufio.NewScanner(_file)
			r := bufio.NewReader(_file)
			if myUrl == false && myDomain == false && myIp == false {
				myUrl = true
			}
			var urlList []string
			var domainList []string
			var ipList []string
			// remove duplicate url
			found := make(map[string]bool)
			//fmt.Println(url, domain, ip)
			//for sc.Scan() {
			for {
				line, err := r.ReadString('\n')
				line = strings.TrimSpace(line)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if err == io.EOF {
					break
				}
				//line := sc.Text()
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
							_domain := searchDomain(_url)
							if _domain == "" {
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
							fmt.Println(_ip)
							found[_ip] = true
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
				for _, line := range urlList {
					_, err := writer.WriteString(line + "\n")
					if err != nil {
						log.Fatal(err)
						return
					}
				}
				for _, line := range domainList {
					_, err := writer.WriteString(line + "\n")
					if err != nil {
						log.Fatal(err)
						return
					}
				}
				for _, line := range ipList {
					_, err := writer.WriteString(line + "\n")
					if err != nil {
						log.Fatal(err)
						return
					}
				}
				err = writer.Flush()
				if err != nil {
					log.Fatal(err)
					return
				}
			}
			//if err := sc.Err(); err != nil {
			//	// line too long occurs error
			//	//log.Fatal(err)
			//	//panic(err)
			//	return
			//}
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
	rootCmd.PersistentFlags().BoolVarP(&myDomain, "domain", "d", false, "search domain from stdin or file")
	rootCmd.PersistentFlags().BoolVarP(&myUrl, "url", "u", false, "search url from stdin or file")
}
