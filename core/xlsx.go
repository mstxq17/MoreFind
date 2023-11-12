package core

import (
	"bytes"
	"github.com/mstxq17/MoreFind/errx"
	"github.com/pbnjay/grate"
	_ "github.com/pbnjay/grate/xls"
	"github.com/tealeg/xlsx"
	"io"
	"os"
	"strings"
)

type OfficeReader interface {
	Read() (io.Reader, error)
}

type XlsxReader struct {
	FilePath string
	BinData  []byte
}

type XlsReader struct {
	FilePath string
	BinData  []byte
}

func NewReader(b []byte) (OfficeReader, error) {
	if len(b) < 4 {
		return nil, errx.NewMsg("Invalid Bytes less then magic number 4 length")
	}
	magicBytes := b[:4]
	var officeReader OfficeReader
	switch {
	case bytes.Equal(magicBytes, []byte{0xD0, 0xCF, 0x11, 0xE0}):
		officeReader = NewXLSReaderFromBinData(b)
	case bytes.Equal(magicBytes, []byte{0x50, 0x4B, 0x03, 0x04}):
		officeReader = NewXLSXReaderFromBinData(b)
	}
	return officeReader, nil
}

// NewXLSXReaderFromFile NewXLSXReader 是一个工厂函数，根据参数创建 xlsxReader 实例
func NewXLSXReaderFromFile(filePath string) *XlsxReader {
	return &XlsxReader{
		FilePath: filePath,
	}
}

// NewXLSXReaderFromBinData 是一个工厂函数，根据参数创建 xlsxReader 实例
func NewXLSXReaderFromBinData(binData []byte) *XlsxReader {
	return &XlsxReader{
		BinData: binData,
	}
}

// NewXLSReaderFromFile  是一个工厂函数，根据参数创建 xlsxReader 实例
func NewXLSReaderFromFile(filePath string) *XlsReader {
	return &XlsReader{
		FilePath: filePath,
	}
}

// NewXLSReaderFromBinData  是一个工厂函数，根据参数创建 xlsxReader 实例
func NewXLSReaderFromBinData(binData []byte) *XlsReader {
	return &XlsReader{
		BinData: binData,
	}
}

func (xr *XlsxReader) Read() (io.Reader, error) {
	var xlsxFile *xlsx.File
	var err error
	if xr.FilePath != "" {
		xlsxFile, err = xlsx.OpenFile(xr.FilePath)
		if err != nil {
			return nil, err
		}

	} else if len(xr.BinData) > 0 {
		xlsxFile, err = xlsx.OpenBinary(xr.BinData)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errx.NewMsg("Value Error, Read xlsx file must provide input")
	}
	var buffer bytes.Buffer
	for _, sheet := range xlsxFile.Sheets {
		for _, row := range sheet.Rows {
			var cellsText []string
			for _, cell := range row.Cells {
				cellsText = append(cellsText, cell.String())
			}
			buffer.WriteString(strings.Join(cellsText, "\t") + NewLine())
		}
	}
	return &buffer, nil
}

func (xr *XlsReader) Read() (io.Reader, error) {
	var wb grate.Source
	var err error
	if xr.FilePath != "" {
		wb, err = grate.Open(xr.FilePath)
		if err != nil {
			return nil, err
		}

	} else if len(xr.BinData) > 0 {
		// create temp file as transfer station
		// 创建一个临时文件做中转站
		tempFile, err := os.CreateTemp("", "temp.*.xls")
		if err != nil {
			return nil, err
		}
		// delete temp file
		// 删除临时文件
		defer os.Remove(tempFile.Name())
		if _, err := tempFile.Write(xr.BinData); err != nil {
			return nil, err
		}
		wb, err = grate.Open(tempFile.Name())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errx.NewMsg("Value Error, Read xls file must provide input")
	}
	var buffer bytes.Buffer
	sheets, _ := wb.List()     // list available sheets
	for _, s := range sheets { // enumerate each sheet name
		sheet, _ := wb.Get(s) // open the sheet
		for sheet.Next() {    // enumerate each row of data
			row := sheet.Strings() // get the row's content as []string
			if len(row) > 0 {
				buffer.WriteString(strings.Join(row, "\t") + NewLine())
			}

		}
	}
	wb.Close()
	return &buffer, nil
}
