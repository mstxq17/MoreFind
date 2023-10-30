package core

import "regexp"

func MatchLine(line, pattern string, inverse bool) (string, error) {
	regexPattern, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	result := regexPattern.MatchString(line)
	if result && !inverse {
		return line, nil
	}
	if !result && inverse {
		return line, nil
	}
	return "", nil
}
