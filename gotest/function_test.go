package gotest

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"log"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

/**
The test standard is that puts the xxx_test.go and xxx.go in same package(directory),
however, some private independent function can be tested in the third party， likes this gotest package。
*/

func isPrivateIP(line string) bool {
	var iIRegex = regexp.MustCompile("^(10.\\d{1,3}.\\d{1,3}.((0/([89]|1[0-9]|2\\d|3[012]))|(\\d{1,3})))|(172.(1[6789]|2\\d|3[01]).\\d{1,3}.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)|(192.168.\\d{1,3}.\\d{1,3}(/(1[6789]|2\\d|3[012]))?)$")
	return iIRegex.MatchString(line)
}

func searchRootDomain(domain string) string {
	eTLD, _ := publicsuffix.EffectiveTLDPlusOne(domain)
	return eTLD
}

func filterExt(fileExt string, filterExts string) bool {
	_exts := strings.Split(filterExts, ",")
	// for improve the filtering speed, reducing the comparative worke，use map
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
	u, err := url.Parse(_url)
	if err != nil {
		log.Fatal(err)
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

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func genIP(cidr string) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("无法解析CIDR地址:", err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		fmt.Println(ip)
	}
}

func Test_genIP(t *testing.T) {
	testCase1 := "192.168.1.0/29"
	genIP(testCase1)
}

func Test_isPrivateIP(t *testing.T) {
	testCase1 := "192.168.1.1"
	testCase2 := "111.210.196.23"
	if isPrivateIP(testCase1) == true && isPrivateIP(testCase2) == false {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}

func Test_searchRootDomain(t *testing.T) {
	testDomains := []string{
		"x.baidu.com",
		"kk.qq.com",
		"x11.xxx.github.io",
		"h.x.中国",
	}
	for _, domain := range testDomains {
		rootDomain := searchRootDomain(domain)
		if domain == rootDomain {
			t.Error("测试失败")
		}
		t.Log(rootDomain + " pass")
	}
	t.Log("全部测试通过")
}

func Test_fileExt(t *testing.T) {
	testUrls := []string{
		"https://baidu.com/",
		"https://baidu.com/123",
		"https://baidu.com/123.png",
		"https://baidu.com/123.png.jpg",
	}
	for index, _url := range testUrls {
		t.Log(strconv.Itoa(index) + ":" + _url)
		switch index {
		case 0:
			t.Log("fileExt:" + fileExt(_url))
			if fileExt(_url) != "" {
				t.Error("测试失败")
			}
		case 1:
			t.Log("fileExt:" + fileExt(_url))
			if fileExt(_url) != "" {
				t.Error("测试失败")
			}
		case 2:
			t.Log("fileExt:" + fileExt(_url))
			if fileExt(_url) != "png" {
				t.Error("测试失败")
			}
		case 3:
			t.Log("fileExt:" + fileExt(_url))
			if fileExt(_url) != "jpg" {
				t.Error("测试失败")
			}
		}
	}
}

func Test_filterExt(t *testing.T) {
	testUrl := "https://baidu.com/123.png"
	if filterExt(fileExt(testUrl), "png, jpg") {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
	testUrl1 := "https://baidu.com/"
	if filterExt(fileExt(testUrl1), "png, jpg") {
		t.Error("测试失败")
	} else {
		t.Log("测试通过")
	}
}
