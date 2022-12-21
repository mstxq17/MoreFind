package cmd

import (
	"net/url"
	"regexp"
	"strings"
)

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
