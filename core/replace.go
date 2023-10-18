package core

import (
	"strings"
)

// Used to parse rules and then replace string contents such as hao123.com -> http://${}/ -> http://123.com/
// 用于解析规则，然后替换字符串内容 比如 hao123.com -> http://${}/ -> http://123.com/

func ReplaceMore(rule string, input string, flag string) string {
	if !strings.Contains(rule, flag) || flag == "" {
		return rule + input
	}
	newOutput := strings.Replace(rule, flag, input, -1)
	return newOutput
}
