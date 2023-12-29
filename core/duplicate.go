package core

import (
	"regexp"
)

const (
	AlphanumericOtherMixed = "{ALPHANUMERIC_OTHER_MIXED}"
	PureNumber             = "{PURE_NUMBER}"
	PureChinese            = "{PURE_CHINESE}"
)

var Filters = map[string]string{
	AlphanumericOtherMixed: `[0-9A-Za-z_-]{8,}`,
	PureNumber:             `[0-9]{2,7}`,
	PureChinese:            `[\p{Han}]{2,}`,
}

// OrderFilters distribute filter order is required because of unordered map
// OrderFilters 组织好过滤顺序是必须的，解决map的无序问题
var OrderFilters = []string{
	AlphanumericOtherMixed,
	PureNumber,
	PureChinese,
}

type DuplicateRemover struct {
	linesMap   map[string]struct{}
	linesCount map[string]int
	threshold  int
	smart      bool
	ANRegexp   map[string]*regexp.Regexp
}

func NewDuplicateRemover(threshold int, smart bool) *DuplicateRemover {
	dr := &DuplicateRemover{
		linesMap:   make(map[string]struct{}),
		linesCount: make(map[string]int),
		threshold:  threshold,
		smart:      smart,
	}
	// some design problems
	// 设计存在问题
	dr.ANRegexp, _ = func() (map[string]*regexp.Regexp, error) {
		ANRegexp := make(map[string]*regexp.Regexp)
		for key, value := range Filters {
			ANRegexp[key] = regexp.MustCompile(value)
		}
		return ANRegexp, nil
	}()
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
	for _, key := range OrderFilters {
		line = dr.ANRegexp[key].ReplaceAllString(line, key)
	}
	return line
}
