package core

import "regexp"

func MatchLine(line, pattern string, inverse bool) (string, error) {
	regexPattern, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	if regexPattern.MatchString(line) && !inverse {
		return line, nil
	}
	if !regexPattern.MatchString(line) && inverse {
		return line, nil
	}
	return "", nil
}
