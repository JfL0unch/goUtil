package excel

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"reflect"
)

type FieldMatchType string

const (
	FieldMatchExactly FieldMatchType = "exactly"
)

type KeyFrom string

const (
	KeyFromTag       KeyFrom = "tag"
	KeyFromFieldName KeyFrom = "field_name"
)

type StructFieldMap map[string]reflect.StructField

type Titles map[int]string

type filter func(e Sheet) Sheet

type Sheet struct {
	titles Titles
	rows   [][]string
}

func (e Sheet) Filter(fs filter) Sheet {
	return fs(e)
}
func (e Sheet) Rows() [][]string {
	return e.rows
}

func (e Sheet) Titles() Titles {
	return e.titles
}
func (e Sheet) Save(fileName string) error {
	excel := excelize.NewFile()

	sheetName := "sheet1"
	err := excel.SetSheetName("Sheet1", sheetName)
	if err != nil {
		return errors.Wrapf(err, "excel.SetSheetName('Sheet1',%s)", sheetName)
	}

	titles := e.Titles()
	err = excel.SetSheetRow(sheetName, fmt.Sprintf("A%d", 1), &titles)
	if err != nil {
		return errors.Wrapf(err, "excel.SetSheetRow(%s,A1)", sheetName)
	}
	for i, row := range e.Rows() {
		cid := i + 2
		err = excel.SetSheetRow(sheetName, fmt.Sprintf("A%d", cid), &row)
		if err != nil {
			fmt.Printf("excel.SetSheetRow(%s,A%d) err:%s", sheetName, cid, err)
		}
	}

	excel.SetActiveSheet(0)

	// Save spreadsheet by the given path.
	if err := excel.SaveAs(fileName); err != nil {
		return err
	}

	return nil
}
