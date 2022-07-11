package gotest

import (
	"golang.org/x/net/publicsuffix"
	"regexp"
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
