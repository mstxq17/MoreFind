package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveDuplicator(t *testing.T) {
	var testCase = []string{
		"http://baidu.com/?a=中国你好呀",
		"http://baidu.com/?a=中国你好呀哈哈",
	}
	dr := NewDuplicateRemover(1, true)
	for _, line := range testCase {
		afterLine := dr.RemoveDuplicator(line)
		require.NotEqualValues(t, "http://baidu.com/?a=中国你好呀哈哈", afterLine)
	}
}

func TestCompareUtf8(t *testing.T) {
	sh := &SimHash{}
	var testCase = []struct {
		s1       string
		s2       string
		expected int64
	}{
		{"w我三个字", "我三个字", 0},
	}
	for _, tc := range testCase {
		result := sh.compareUtf8(tc.s1, tc.s2)
		require.Equal(t, tc.expected, result, "测试失败")
	}
}
