package excel

type Intf interface {
	GetSheet(excelFile, sheetName string, opts ...Opt) (*Sheet, error)
}
