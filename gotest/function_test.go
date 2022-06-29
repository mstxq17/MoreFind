package gotest

import (
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

func Test_isPrivateIP(t *testing.T) {
	testCase1 := "192.168.1.1"
	testCase2 := "111.210.196.23"
	if isPrivateIP(testCase1) == true && isPrivateIP(testCase2) == false {
		t.Log("测试通过")
	} else {
		t.Error("测试失败")
	}
}
