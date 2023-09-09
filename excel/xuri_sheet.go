package excel

import (
	"github.com/xuri/excelize/v2"
)

type Xuri struct{}

func (x Xuri) GetSheet(excelFile, sheetName string, opts ...Opt) (*Sheet, error) {

	parser := &Parser{}

	for _, opt := range opts {
		opt(parser)
	}
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		return nil, err
	}

	excelDatas, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(excelDatas) <= 0 {
		return nil, nil
	}

	excel := Sheet{}
	titles, err := parser.GetTitles(excelDatas)
	if err != nil {
		return nil, err
	}
	excel.titles = titles

	rows, err := parser.GetRows(excelDatas)
	if err != nil {
		return nil, err
	}
	excel.rows = rows

	return &excel, nil
}
