package main

import (
	"github.com/mstxq17/MoreFind/cmd"
)

/**
程序执行流程如下
1）解析参数 -domain -url -ip
2）从管道读取输入
3）输出结果
*/

func main() {
	cmd.Execute()
}
