package core

import (
	"fmt"
	"os"
)

func NewLine() string {
	var PS = fmt.Sprintf("%v", os.PathSeparator)
	var LineBreak = "\n"
	if PS != "/" {
		LineBreak = "\r\n"
	}
	return LineBreak
}
