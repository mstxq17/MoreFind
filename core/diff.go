package core

import (
	"bufio"
	"os"
	"sort"
)

// 用于比较两个文件a、b每一行,根据需要提取以下情况
// 1: a有的，b没有的行
// 2: a没有的，b有的行
// 3: a、b都有的行
// 严格模式是逐行比较
// 非严格模式是排序后逐行比较，默认

func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func CompareFiles(a, b []string, strictMode bool) ([]string, []string, []string) {
	var onlyInA, onlyInB, inBoth []string

	if strictMode {
		// 严格模式：逐行比较
		lineACount := len(a)
		lineBCount := len(b)
		for index, lineA := range a {
			if index <= lineBCount-1 {
				if lineA != b[index] {
					onlyInA = append(onlyInA, lineA)
				} else {
					if lineACount > lineBCount {
						inBoth = append(inBoth, lineA)
					}
				}
			} else {
				onlyInA = append(onlyInA, lineA)
			}
		}
		for index, lineB := range b {
			if index <= lineACount-1 {
				if lineB != a[index] {
					onlyInB = append(onlyInA, lineB)
				} else {
					if lineBCount >= lineACount {
						inBoth = append(inBoth, lineB)
					}
				}
			} else {
				onlyInB = append(onlyInB, lineB)
			}
		}
	} else {
		// 非严格模式：排序后逐行比较
		sort.Strings(a)
		sort.Strings(b)

		// 逐行比较
		ia, ib := 0, 0
		for ia < len(a) && ib < len(b) {
			switch {
			case a[ia] == b[ib]:
				inBoth = append(inBoth, a[ia])
				ia++
				ib++
			case a[ia] < b[ib]:
				onlyInA = append(onlyInA, a[ia])
				ia++
			default:
				onlyInB = append(onlyInB, b[ib])
				ib++
			}
		}

		// 处理剩余部分
		for ; ia < len(a); ia++ {
			onlyInA = append(onlyInA, a[ia])
		}
		for ; ib < len(b); ib++ {
			onlyInB = append(onlyInB, b[ib])
		}
	}

	return onlyInA, onlyInB, inBoth
}
