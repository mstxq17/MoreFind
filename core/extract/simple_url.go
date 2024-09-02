package extract

import (
	"fmt"
	"net/url"
)

func SimpleUrl(input string) (string, bool, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", false, err
	}
	if parsed.Scheme != "" {
		tSimpleUrl := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
		return tSimpleUrl, true, nil
	} else {
		//  拼接协议进行解析
		modifiedHost := fmt.Sprintf("http://%s", input)
		parsed, err := url.Parse(modifiedHost)
		if err != nil {
			return "", false, err
		} else {
			tSimpleUrl := fmt.Sprintf("%s", parsed.Hostname())
			return tSimpleUrl, false, nil
		}
	}
}
