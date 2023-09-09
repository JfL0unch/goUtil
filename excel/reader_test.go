package excel

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type Status uint8
type typX struct {
	Id     uint64    `json:"id"`
	Name   string    `json:"name"`
	Point  float64   `json:"point"`
	Time   time.Time `json:"time"`
	Status Status    `json:"status"`
}

func TestReader_getStructFieldMap(t *testing.T) {

	FocusConvey("1", t, func() {
		Convey("no Anonymous", func() {
			config := ReaderConfig{
				SheetWithTitle: true,
				KeyFrom:        KeyFromTag,
				KeyTagName:     "json",
			}
			p := NewReader(config)

			x := typX{}
			mp, err := p.getStructFieldMap(x)
			So(err, ShouldBeNil)

			So(mp["id"].Name, ShouldEqual, "Id")
			So(mp["name"].Name, ShouldEqual, "Name")
			So(mp["time"].Name, ShouldEqual, "Time")
			So(mp["status"].Name, ShouldEqual, "Status")

		})

		FocusConvey("Anonymous", func() {
			config := ReaderConfig{
				SheetWithTitle: true,
				KeyFrom:        KeyFromTag,
				KeyTagName:     "json",
			}
			p := NewReader(config)

			x := typX{}
			mp, err := p.getStructFieldMap(x)
			So(err, ShouldBeNil)

			So(mp["id"].Name, ShouldEqual, "Id")
			So(mp["name"].Name, ShouldEqual, "Name")

		})
	})
}
func TestReader_Parse(t *testing.T) {
	FocusConvey("0", t, func() {

		Convey("data", func() {
			config := ReaderConfig{
				SheetWithTitle: true,
				KeyFrom:        KeyFromTag,
				KeyTagName:     "json",
			}
			p := NewReader(config)

			x := typX{}

			excelFile := "data.xlsx"
			retI, err := p.Parse(x, excelFile, "Sheet1")

			So(p.sheet, ShouldNotBeNil)
			So(len(p.sheet.Rows()), ShouldEqual, 4)
			So(p.structFieldMap, ShouldNotBeNil)
			So(p.structFieldMap["id"].Name, ShouldEqual, "Id")
			So(p.structFieldMap["name"].Name, ShouldEqual, "Name")
			So(p.structFieldMap["time"].Name, ShouldEqual, "Time")
			So(p.structFieldMap["status"].Name, ShouldEqual, "Status")

			ret, ok := retI.([]typX)
			So(ok, ShouldEqual, true)
			So(err, ShouldBeNil)

			So(len(ret), ShouldEqual, 4)
			So(ret[0].Id, ShouldEqual, 1)
			So(ret[0].Name, ShouldEqual, "jack")
			So(ret[0].Point, ShouldEqual, 17.23)
			So(ret[0].Time, ShouldEqual, time.Date(2023, time.August, 7, 0, 34, 0, 0, time.Local))
			So(ret[0].Status, ShouldEqual, 1)

			So(ret[1].Id, ShouldEqual, 2)
			So(ret[1].Name, ShouldEqual, "tom")
			So(ret[1].Point, ShouldEqual, 0.58)
			So(ret[1].Time, ShouldEqual, time.Date(2023, time.August, 7, 0, 34, 51, 0, time.Local))
			So(ret[1].Status, ShouldEqual, 1)

			So(ret[2].Id, ShouldEqual, 3)
			So(ret[2].Name, ShouldEqual, "lucy")
			So(ret[2].Point, ShouldEqual, -3)
			So(ret[2].Time, ShouldEqual, time.Date(2023, time.August, 26, 7, 40, 31, 0, time.Local))
			So(ret[2].Status, ShouldEqual, 2)

			So(ret[3].Id, ShouldEqual, 4)
			So(ret[3].Name, ShouldEqual, "kitty")
			So(ret[3].Point, ShouldEqual, -1.1)
			So(ret[3].Time, ShouldEqual, time.Date(2023, time.August, 25, 0, 0, 0, 0, time.Local))
			So(ret[3].Status, ShouldEqual, 2)
		})
	})
}
