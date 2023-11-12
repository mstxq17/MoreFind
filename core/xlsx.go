package core

import (
	"github.com/tealeg/xlsx"
	"io"
	"strings"
)

type OfficeReader interface {
	ReadXLSX() (*xlsx.File, error)
}

type FileReader struct {
	FilePath string
}

type BinReader struct {
	BinData []byte
}

func (fr *FileReader) ReadXLSX() (*xlsx.File, error) {
	xlsxFile, err := xlsx.OpenFile(fr.FilePath)
	if err != nil {
		return nil, err
	}
	return xlsxFile, nil
}

func (br *BinReader) ReadXLSX() (*xlsx.File, error) {
	xlsxFile, err := xlsx.OpenBinary(br.BinData)
	if err != nil {
		return nil, err
	}
	return xlsxFile, nil
}

func ReadIOReader(xf *xlsx.File) (io.Reader, error) {
	var allRowTexts []string
	// work through sheets
	// 遍历每一个工作表
	for _, sheet := range xf.Sheets {
		// 遍历每一行
		for _, row := range sheet.Rows {
			// 存储单行的字符串
			var cellsText []string
			// 遍历每个单元格
			for _, cell := range row.Cells {
				text := cell.String()
				cellsText = append(cellsText, text)
			}
			rowText := strings.Join(cellsText, ",")
			allRowTexts = append(allRowTexts, rowText)
		}
	}
	allRowTextStrings := strings.Join(allRowTexts, "\n")
	return strings.NewReader(allRowTextStrings), nil
}
