package core

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestGenIP(t *testing.T) {
	testCase := "192.168.1.1-192.168.1.2"
	expectOut := 2
	outputchan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		count := 0
		for o := range outputchan {
			fmt.Println(o)
			count++
		}
		require.Equalf(t, expectOut, count, "输出的结果总数不对")
		wg.Done()
	}()
	err := GenIP(testCase, outputchan)
	require.Nil(t, err)
	close(outputchan)
	wg.Wait()
}
