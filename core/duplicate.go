package core

import "regexp"

const (
	threshold = 10
)
const (
	ALPHANUMERIC = "{ALPHANUMERIC}"
)

type DuplicateRemover struct {
	linesMap   map[string]struct{}
	linesCount map[string]int
	threshold  int
	smart      bool
	ANRegexp   *regexp.Regexp
}

func NewDuplicateRemover(threshold int, smart bool) *DuplicateRemover {
	dr := &DuplicateRemover{
		linesMap:   make(map[string]struct{}),
		linesCount: make(map[string]int),
		threshold:  threshold,
		smart:      smart,
	}
	dr.ANRegexp, _ = regexp.Compile(`[0-9A-Za-z]{10,}`)
	return dr
}

func (dr *DuplicateRemover) RemoveDuplicator(line string) string {
	if dr.smart {
		gResult := dr.generalize(line)
		if _, exists := dr.linesMap[line]; !exists {
			dr.linesMap[line] = struct{}{}
			dr.linesCount[gResult] += 1
			if dr.linesCount[gResult] <= dr.threshold {
				return line
			}
		}
	} else {
		if _, exists := dr.linesMap[line]; !exists {
			dr.linesMap[line] = struct{}{}
			return line
		}
	}
	return ""
}

// 将正则 [0-9A-Za-z]{10,} 一般化，超过阈值则进行智能过滤
func (dr *DuplicateRemover) generalize(line string) string {
	return dr.ANRegexp.ReplaceAllString(line, ALPHANUMERIC)
}
