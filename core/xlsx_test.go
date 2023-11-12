package core

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestOfficeReader(t *testing.T) {
	testCases := []struct {
		input  string
		expect interface{}
	}{
		{input: "Sample.xlsx", expect: nil},
		{input: "Sample.xls", expect: nil},
	}
	for _, testCase := range testCases {
		input := testCase.input
		expect := testCase.expect
		if strings.HasSuffix(input, "xlsx") {
			var reader OfficeReader = NewXLSXReaderFromFile(input)
			_, err := reader.Read()
			//fmt.Println(content)
			require.Equal(t, expect, err)
			binData, _ := os.ReadFile(input)
			reader = NewXLSXReaderFromBinData(binData)
			content, err := reader.Read()
			fmt.Println(content)
			require.Equal(t, expect, err)
			reader, err = NewReader(binData)
			require.Equal(t, expect, err)
		}
		if strings.HasSuffix(input, "xls") {
			var reader OfficeReader = NewXLSReaderFromFile(input)
			_, err := reader.Read()
			require.Equal(t, expect, err)
			binData, _ := os.ReadFile(input)
			reader = NewXLSReaderFromBinData(binData)
			_, err = reader.Read()
			//fmt.Println(content)
			require.Equal(t, expect, err)
			reader, err = NewReader(binData)
			require.Equal(t, expect, err)
		}
	}
}
