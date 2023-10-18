package core

import (
	"bytes"
)

type stringHandler interface {
	HandleString(s string) string
}

type MyBuffer struct {
	buffer   *bytes.Buffer
	IsFilter bool
}

type CustomStringHandler struct {
	Rule     string
	Flag     string
	Strategy int //类型
}

func (c *CustomStringHandler) HandleString(s string) string {
	if c.Strategy == 1 {
		return ReplaceMore(c.Rule, s, c.Flag)
	} else {
		return s
	}

}

func NewMyBuffer(isFilter bool) *MyBuffer {
	return &MyBuffer{
		buffer:   new(bytes.Buffer),
		IsFilter: isFilter,
	}
}

func (_bytes *MyBuffer) WriteString(s string, handler stringHandler, newLine string) (n int, err error) {
	// change the action of WriteString method
	// 修改 WriteString 方法的行为
	if _bytes.IsFilter {
		return _bytes.buffer.WriteString(handler.HandleString(s) + newLine)
	}
	return _bytes.buffer.WriteString(s + newLine)
}

func (_bytes *MyBuffer) String() string {
	return _bytes.buffer.String()
}

func (_bytes *MyBuffer) Reset() {
	_bytes.buffer.Reset()
}
