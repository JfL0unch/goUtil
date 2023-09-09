package excel

import (
	"fmt"
	"github.com/JfL0unch/goUtil/reflectUtils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

type ReaderConfig struct {
	// true:
	//		column name  of sheet  to field name of struct.
	//
	// false:
	//		column index of sheet  to field index of struct,
	//		one to one strictly.
	// 		len(field) < len(column).
	SheetWithTitle bool

	KeyFrom    KeyFrom
	KeyTagName string // todo 支持 gorm.column这种格式
}
type Reader struct {
	config         ReaderConfig
	structTmpl     interface{}
	sheet          *Sheet
	structFieldMap StructFieldMap
}

func NewReader(c ReaderConfig) Reader {
	return Reader{
		config: c,
	}
}

func (r *Reader) paramCheckOk(propType interface{}, excelFile, sheetName string) (msg string, ok bool) {
	rType := reflect.TypeOf(propType)
	rKind := rType.Kind()

	// 1. propType should be sliced
	if rKind != reflect.Struct {
		return fmt.Sprintf("propType(%s) not struct", rKind.String()), false
	}

	// todo excelFile-sheet should be accessible
	return "", true
}

func (r *Reader) Parse(structTmpl interface{}, excelFile, sheetName string) (interface{}, error) {
	if msg, ok := r.paramCheckOk(structTmpl, excelFile, sheetName); !ok {
		return nil, errors.New("r.paramCheckOk failed:" + msg)
	}

	r.structTmpl = structTmpl

	x := Xuri{}
	sheet, err := x.GetSheet(excelFile, sheetName, FirstRowAsTitles())
	if err != nil {
		return nil, errors.Wrapf(err, "GetSheet(%s,%s)", excelFile, sheetName)
	}

	r.sheet = sheet

	structFieldMap, err := r.getStructFieldMap(structTmpl)
	if err != nil {
		return nil, errors.Wrapf(err, "getStructFieldMap(%v)", structTmpl)
	}
	r.structFieldMap = structFieldMap

	return r.ProcessRows()

}

func (r *Reader) getStructInstance(columns []string) reflect.Value {

	fieldMap := r.structFieldMap
	sheet := r.sheet
	sheetTitles := sheet.Titles()
	sheetWithTitle := r.config.SheetWithTitle

	// 1
	structProto := r.structTmpl
	structTyp := reflect.TypeOf(structProto)
	structInstance := reflect.New(structTyp)

	if sheetWithTitle {
		for columnIndex, columnStr := range columns {
			fieldName := sheetTitles[columnIndex]
			parseWithTitle(structInstance, fieldName, columnStr, fieldMap)
		}
	} else {
		for columnIndex, columnVal := range columns {
			fieldTmpl := structInstance.Field(columnIndex)
			parsedVal, err := reflectUtils.ParseStrToInstance(fieldTmpl, columnVal)
			if err == nil {
				fieldTmpl.Set(parsedVal)
			}
		}
	}

	return structInstance.Elem()
}

func (r *Reader) ProcessRows() (interface{}, error) {

	sheet := r.sheet
	propStruct := r.structTmpl

	structTyp := reflect.TypeOf(propStruct)

	capSize := len(sheet.rows)
	structSlice := reflect.MakeSlice(reflect.SliceOf(structTyp), 0, capSize)

	for _, row := range sheet.rows {

		structInstance := r.getStructInstance(row)

		structSlice = reflect.Append(structSlice, structInstance)
	}

	return structSlice.Interface(), nil

}

func (r *Reader) getStructFieldMap(structTmpl interface{}) (StructFieldMap, error) {
	res := make(StructFieldMap, 0)
	ift := reflect.TypeOf(structTmpl)
	ifv := reflect.ValueOf(structTmpl)
	keyFrom := r.config.KeyFrom
	tagName := r.config.KeyTagName
	if tagName == "" {
		tagName = "json"
	}
	var key func(_ reflect.StructField) string

	switch keyFrom {
	case KeyFromTag:
		key = func(field reflect.StructField) string {
			return field.Tag.Get(tagName)
		}
	case KeyFromFieldName:
		fallthrough
	default:
		key = func(field reflect.StructField) string {
			return field.Name
		}
	}

	for i := 0; i < ifv.NumField(); i++ {
		ft := ift.Field(i)
		fv := ifv.Field(i)

		if ft.Type.Kind() == reflect.Struct && ft.Anonymous {
			deepFields := reflectUtils.FlatStructFields(fv.Interface())
			if len(deepFields) > 0 {
				for _, df := range deepFields {
					keyName := key(df)
					res[keyName] = df
				}
			}
		} else {
			keyName := key(ft)

			res[keyName] = ft
		}

	}
	return res, nil
}

type Opt func(p *Parser)

func FirstRowAsTitles() Opt {
	return func(p *Parser) {
		p.withTitles = true
	}
}

type Parser struct {
	// 第一行是否列名
	withTitles bool
}

func (p Parser) WithTitle() bool {
	return p.withTitles
}

func (p Parser) GetTitles(excelData [][]string) (Titles, error) {
	if len(excelData) <= 0 {
		return nil, errors.New("GetTitles err: no data")
	}

	titles := make(Titles, 0)
	for i, v := range excelData[0] {
		if p.WithTitle() {
			titles[i] = v
		} else {
			titles[i] = strconv.FormatInt(int64(i), 10)
		}
	}
	return titles, nil
}

func (p Parser) GetRows(excelData [][]string) ([][]string, error) {
	if len(excelData) <= 0 {
		return nil, nil
	}
	if p.WithTitle() {
		return excelData[1:], nil
	} else {
		return excelData, nil
	}
}

func parseWithTitle(structToUpdate reflect.Value, fieldName, columnVal string, structMap StructFieldMap) {
	if field, existField := structMap[fieldName]; existField {
		fieldTmpl := structToUpdate.Elem().FieldByName(field.Name)
		nv, err := reflectUtils.ParseStrToInstance(fieldTmpl, columnVal)
		if err == nil {
			fieldTmpl.Set(nv)
		} else {
			logrus.Errorf("reflectUtils.ParseStrToInstance err:%s", err)
		}
	}
}
