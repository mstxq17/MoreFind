package core

import (
	"fmt"
	"unicode/utf8"
)

// Universal Similarity Recognition
// 通用相似度识别

// simhash implementation
// SimHash 算法实现

type SimHash struct {
}

// 分词&权重
func (sh *SimHash) compareUtf8(s1, s2 string) float64 {
	s1Len := utf8.RuneCountInString(s1)
	s2Len := utf8.RuneCountInString(s2)
	fmt.Println(s1Len, s2Len, len(s1), len(s2))
	return 0
}
