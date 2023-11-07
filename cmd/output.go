package cmd

import (
	"fmt"
	"github.com/mstxq17/MoreFind/core"
	"github.com/mstxq17/MoreFind/errx"
	"os"
	"sync"
)

var NewLine = core.NewLine()

// try for refactor in future
// 为以后统一输出做铺垫
func syncOutput(wg *sync.WaitGroup, outputchan chan string) {
	// 任务完成,增加计数
	defer wg.Done()
	var f *os.File
	if output != "" {
		var err error
		f, err = os.Create(output)
		if err != nil {
			logger.Fatal(errx.NewWithMsgf(err, "Could not create output file '%s':", file))
		}
		defer f.Close()
	}
	for o := range outputchan {
		if o == "" {
			continue
		}
		// output to stdout & file stream
		// 输出到 stdout & 文件流
		if len(myIPFormats) > 0 {
			outputItems(f, core.AlterIP(o, myIPFormats)...)
		} else {
			outputItems(f, o)
		}
	}
}

func outputItems(f *os.File, items ...string) {
	for _, item := range items {
		fmt.Println(item)
		if f != nil {
			_, _ = f.WriteString(item + NewLine)
		}
	}
}
