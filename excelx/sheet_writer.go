package excelx

import (
	"github.com/tealeg/xlsx"
	
	"github.com/go-xuan/utilx/errorx"
)

type SheetWriter interface {
	WriteHeader(sheet *xlsx.Sheet)
	WriteRow(sheet *xlsx.Sheet)
}

// WriteSheet 将数据写入sheet页
func WriteSheet[W SheetWriter](file *xlsx.File, sheetName string, data []W) error {
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return errorx.Wrap(err, "add sheet error")
	}
	if len(data) > 0 {
		data[0].WriteHeader(sheet)
		for _, row := range data {
			row.WriteRow(sheet)
		}
	}
	return nil
}
