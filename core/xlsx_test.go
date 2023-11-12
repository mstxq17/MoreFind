package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOfficeReader(t *testing.T) {
	testCaseFile := "Sample.xlsx"
	fileReader := &FileReader{FilePath: testCaseFile}
	fr, err := fileReader.ReadXLSX()
	require.Nil(t, err)
	_, err = ReadIOReader(fr)
	require.Nil(t, err)
}
