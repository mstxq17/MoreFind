package extract

import (
	"fmt"
	"net/url"
	"strings"
)

func SimpleUrl(input string) (string, bool, error) {
	var hasSchema bool
	if strings.Contains(input, "://") {
		hasSchema = true
	} else {
		hasSchema = false
		input = fmt.Sprintf("http://%s", input)
	}
	parsed, err := url.Parse(input)
	if err != nil {
		return "", false, err
	}
	if hasSchema == true {
		tSimpleUrl := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
		return tSimpleUrl, true, nil
	} else {
		//  拼接协议进行解析
		if err != nil {
			return "", false, err
		} else {
			tSimpleUrl := fmt.Sprintf("%s", parsed.Host)
			return tSimpleUrl, false, nil
		}
	}
}
