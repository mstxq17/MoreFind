package gotest

import (
	"net/url"
	"regexp"
	"strings"
	"testing"
)

type URL struct {
	Scheme string // protocol
	Host   string // host or  host:port
	Path   string // relative paths may omit leading slash
}

func deepMakeUrl(checkUrl string) string {
	_, err := url.PathUnescape(checkUrl)
	if err != nil {
		errStr := err.Error()
		re := regexp.MustCompile(`invalid URL escape "(.+)"`)
		subString := re.FindStringSubmatch(errStr)[1]
		idx := strings.LastIndex(checkUrl, subString)
		if idx >= 0 {
			checkUrl = checkUrl[:idx] + "" + checkUrl[idx+len(subString):]
		}
		return deepMakeUrl(checkUrl)
	}
	return checkUrl
}

func Test_deepMakeUrl(t *testing.T) {
	testCases := []struct {
		url       string
		excpected string
	}{
		{"https://baidu.com/%23", "https://baidu.com/%23"},
		{"https://baidu.com/%23%1", "https://baidu.com/%23"},
		{"https://baidu.com/%23%", "https://baidu.com/%23"},
		{"https://baidu.com/%23%1%1", "https://baidu.com/%231"},
		{"https://baidu.com/%232%%", "https://baidu.com/%232"},
		{"https://baidu.com/%232%%23", "https://baidu.com/%2323"},
	}
	for _, tc := range testCases {
		_, err := url.Parse(tc.excpected)
		if err != nil {
			t.Errorf("Expected %s not valid", tc.excpected)
		}
		result := deepMakeUrl(tc.url)
		if result != tc.excpected {
			t.Errorf("Origin %s Expected %s, got %s", tc.url, tc.excpected, result)
		}
	}
}
