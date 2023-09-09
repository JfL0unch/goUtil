package reflectUtils

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"time"
)

func Test_isAlias(t *testing.T) {

	Convey("1", t, func() {

		Convey("string not alias", func() {
			var i string
			v := reflect.ValueOf(i)

			got := isAliasType(v)
			So(got, ShouldEqual, false)
		})

		Convey("string alias", func() {
			type strAlias string
			var i strAlias
			v := reflect.ValueOf(i)

			got := isAliasType(v)
			So(got, ShouldEqual, true)
		})

		Convey("uint8 alias", func() {
			type int8Alias uint8
			var i int8Alias
			v := reflect.ValueOf(i)

			got := isAliasType(v)
			So(got, ShouldEqual, true)
		})

	})
}

func Test_convertAliasType(t *testing.T) {

	Convey("1", t, func() {

		Convey("string not alias", func() {
			strVal := "8"
			var val uint8
			valI := reflect.ValueOf(val)
			instance, err := getInstance(valI, strVal)

			gotI, err := getInstanceOfAliasType(instance, strVal)
			So(err, ShouldBeNil)
			got := gotI.Interface()
			_, assertOk := got.(uint8)
			So(assertOk, ShouldEqual, true)
			So(got, ShouldEqual, 8)
		})

		Convey("uint8 alias", func() {
			strVal := "8"
			type int8Alias uint8

			var i int8Alias
			v := reflect.ValueOf(i)
			gotI, err := getInstanceOfAliasType(v, strVal)
			So(err, ShouldBeNil)
			got := gotI.Interface()
			_, assertOk := got.(int8Alias)
			So(assertOk, ShouldEqual, true)
			So(got, ShouldEqual, int8Alias(8))
		})

	})
}

func Test_ParseStrToTmpl(t *testing.T) {

	Convey("1", t, func() {
		Convey("string", func() {
			strVal := "test"

			var i string
			v := reflect.ValueOf(i)
			got, err := ParseStrToInstance(v, strVal)

			So(err, ShouldBeNil)
			So(got.Interface(), ShouldEqual, "test")
			So(got.Interface(), ShouldNotEqual, "123")
		})

		Convey("string alias", func() {
			strVal := "test"

			type strAlias string
			var i strAlias
			v := reflect.ValueOf(i)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)
			got := gotI.Interface()
			_, assertOk := got.(strAlias)
			So(assertOk, ShouldEqual, true)
			So(got, ShouldEqual, "test")
			So(got, ShouldNotEqual, "123")
		})

		Convey("int", func() {
			strVal := "-13"

			var i int
			v := reflect.ValueOf(i)
			got, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)
			So(got.Interface(), ShouldEqual, -13)
			So(got.Interface(), ShouldNotEqual, 13)
		})

		Convey("bool", func() {
			strVal := "false"
			var i bool

			v := reflect.ValueOf(i)
			got, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			So(got.Interface(), ShouldEqual, false)
			So(got.Interface(), ShouldNotEqual, true)
		})
		Convey("uint", func() {
			strVal := "64301"

			var i uint
			v := reflect.ValueOf(i)
			got, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			So(got.Interface(), ShouldEqual, 64301)
			So(got.Interface(), ShouldNotEqual, 65536)
		})

		Convey("struct", func() {
			strVal := `{"name":"tom","age":13,"isMan":true}`

			type typX struct {
				Name  string `json:"name"`
				Age   int    `json:"age"`
				IsMan bool   `json:"isMan"`
			}
			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got.Name, ShouldEqual, "tom")
			So(got.Age, ShouldEqual, 13)
			So(got.IsMan, ShouldEqual, true)
		})

		FocusConvey("struct no-datetime alias", func() {
			strVal := `{"name":"tom","age":13,"isMan":true,"status":1}`

			type Status uint8
			type typX struct {
				Name   string `json:"name"`
				Age    int    `json:"age"`
				IsMan  bool   `json:"isMan"`
				Status Status `json:"status"`
			}

			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got.Name, ShouldEqual, "tom")
			So(got.Age, ShouldEqual, 13)
			So(got.IsMan, ShouldEqual, true)
			So(got.Status, ShouldEqual, 1)
		})

		Convey("struct datetime alias", func() {
			strVal := `{"name":"tom","age":13,"isMan":true,"date":"2023-08-13 00:00:00","hpa_status":1}`

			type Status uint8
			type typX struct {
				Name      string    `json:"name"`
				Age       int       `json:"age"`
				IsMan     bool      `json:"isMan"`
				Date      time.Time `json:"date"`
				HpaStatus Status    `json:"hpa_status"`
			}

			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got.Name, ShouldEqual, "tom")
			So(got.Age, ShouldEqual, 13)
			So(got.IsMan, ShouldEqual, true)
			So(got.HpaStatus, ShouldEqual, 1)
		})

		Convey("struct embedded", func() {
			strVal := `{"name":"tom","age":13,"info":{"sex":1,"height":128.3,"friend":"jack"}}`

			type friend string
			type info struct {
				Sex    int     `json:"sex"`
				Height float64 `json:"height"`
				Friend friend  `json:"friend"`
			}
			type typX struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Info info   `json:"info"`
			}
			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got.Name, ShouldEqual, "tom")
			So(got.Age, ShouldEqual, 13)
			So(got.Info.Height, ShouldEqual, 128.3)
			So(got.Info.Sex, ShouldEqual, 1)
			So(got.Info.Friend, ShouldEqual, "jack")

		})
		Convey("map", func() {
			strVal := `{"name":"tom","age":13,"height":121.93,"isMan":true}`

			type typX map[string]interface{}
			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got["name"], ShouldEqual, "tom")
			So(got["age"], ShouldEqual, 13)
			So(got["height"], ShouldEqual, 121.93)
			So(got["isMan"], ShouldEqual, true)

		})

		Convey("array int", func() {
			strVal := `[1.0,3,11.1]`

			type typX [3]float64
			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got[0], ShouldEqual, 1.0)
			So(got[1], ShouldEqual, 3)
			So(got[2], ShouldEqual, 11.1)
		})
		Convey("array struct", func() {
			strVal := `[{"name":"tom","age":13,"height":121.93,"isMan":true},{"name":"jack","age":11,"height":98.93,"isMan":false}]`

			type typX struct {
				Name   string  `json:"name"`
				Age    int     `json:"age"`
				Height float64 `json:"height"`
				IsMan  bool    `json:"isMan"`
			}
			v := reflect.ValueOf([2]typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			gots := gotI.Interface().([2]typX)
			got := gots[0]
			So(got.Name, ShouldEqual, "tom")
			So(got.Age, ShouldEqual, 13)
			So(got.Height, ShouldEqual, 121.93)
			So(got.IsMan, ShouldEqual, true)

			got1 := gots[1]
			So(got1.Name, ShouldEqual, "jack")
			So(got1.Age, ShouldEqual, 11)
			So(got1.Height, ShouldEqual, 98.93)
			So(got1.IsMan, ShouldEqual, false)
		})

		Convey("slice int", func() {
			strVal := `[1.0,3,11.1]`

			type typX []float64
			v := reflect.ValueOf(typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(typX)

			So(got[0], ShouldEqual, 1.0)
			So(got[1], ShouldEqual, 3)
			So(got[2], ShouldEqual, 11.1)

		})

		Convey("slice struct", func() {
			strVal := `[{"name":"tom","age":13,"height":121.93,"isMan":true},{"name":"jack","age":11,"height":98.93,"isMan":false}]`

			type typX struct {
				Name   string  `json:"name"`
				Age    int     `json:"age"`
				Height float64 `json:"height"`
				IsMan  bool    `json:"isMan"`
			}
			v := reflect.ValueOf([]typX{})
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			gots := gotI.Interface().([]typX)
			got := gots[0]
			So(got.Name, ShouldEqual, "tom")
			So(got.Age, ShouldEqual, 13)
			So(got.Height, ShouldEqual, 121.93)
			So(got.IsMan, ShouldEqual, true)

			got1 := gots[1]
			So(got1.Name, ShouldEqual, "jack")
			So(got1.Age, ShouldEqual, 11)
			So(got1.Height, ShouldEqual, 98.93)
			So(got1.IsMan, ShouldEqual, false)
		})

		Convey("pointer int", func() {
			strVal := `98`

			var typeX int
			typeY := &typeX
			v := reflect.ValueOf(typeY)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*int)

			So(*got, ShouldEqual, 98)

		})

		Convey("pointer with input string('')", func() {
			strVal := ``

			var typeX struct{}

			typeY := &typeX
			v := reflect.ValueOf(typeY)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*struct{})

			So(*got, ShouldResemble, typeX)

		})

		Convey("struct pointer ", func() {
			strVal := `{"age":13} `

			type typeY struct {
				Age *int `json:"age"`
			}

			typeX := reflect.ValueOf(&typeY{})
			gotI, err := ParseStrToInstance(typeX, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*typeY)

			So(got, ShouldNotBeNil)

			So(*got.Age, ShouldEqual, 13)

		})

		Convey("struct pointer with input string('')", func() {
			strVal := ""

			var typeX *struct {
				Age *int `json:"age"`
			}
			v := reflect.ValueOf(typeX)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*struct {
				Age *int `json:"age"`
			})

			So(got, ShouldBeNil)

		})

		Convey("struct with input string(null) ", func() {
			strVal := "null"

			var typeX *struct {
				Age *int `json:"age"`
			}
			v := reflect.ValueOf(typeX)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*struct {
				Age *int `json:"age"`
			})

			So(got, ShouldBeNil)

		})

		Convey("time.Time(2023-08-25 14:30:31) ", func() {
			strVal := "2023-08-25 14:30:31"

			var typeX time.Time
			v := reflect.ValueOf(&typeX)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*time.Time)

			So(*got, ShouldEqual, time.Date(2023, time.August, 25, 14, 30, 31, 0, time.Local))

		})
		Convey("time.Time(8/21/18 23:38) ", func() {
			strVal := "8/21/18 23:38"

			var typeX time.Time
			v := reflect.ValueOf(&typeX)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*time.Time)

			So(*got, ShouldEqual, time.Date(2018, time.August, 21, 23, 38, 0, 0, time.Local))

		})
		Convey("time.Time(4/22/22 16:24) ", func() {
			strVal := "4/22/22 16:24"

			var typeX time.Time
			v := reflect.ValueOf(&typeX)
			gotI, err := ParseStrToInstance(v, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*time.Time)

			So(*got, ShouldEqual, time.Date(2022, time.April, 22, 16, 24, 0, 0, time.Local))

		})

		Convey("struct float64(E+12)", func() {
			strVal := `{"length":2.98109E+12} `

			type typeY struct {
				Length float64 `json:"length"`
			}

			typeX := reflect.ValueOf(&typeY{})
			gotI, err := ParseStrToInstance(typeX, strVal)
			So(err, ShouldBeNil)

			got := gotI.Interface().(*typeY)

			So(got, ShouldNotBeNil)

			So(got.Length, ShouldEqual, 2.98109e+12)

		})

	})

}

func Test_getTimeFromStr(t *testing.T) {

	Convey("1", t, func() {
		Convey("2023-08-13 00:00:00", func() {
			got, err := getTimeFromStr("2023-08-13 00:00:00")
			So(err, ShouldBeNil)
			So(got, ShouldEqual, time.Date(2023, time.August, 13, 0, 0, 0, 0, time.Local))
		})
		Convey("2023-07-04T09:36:33.961605775", func() {
			got, err := getTimeFromStr("2023-07-04T09:36:33.961605775")
			So(err, ShouldBeNil)
			So(got, ShouldEqual, time.Date(2023, time.July, 4, 9, 36, 33, 961605775, time.Local))
		})

		Convey("2023-07-04T09:36:33", func() {
			got, err := getTimeFromStr("2023-07-04T09:36:33")
			So(err, ShouldBeNil)
			So(got, ShouldEqual, time.Date(2023, time.July, 4, 9, 36, 33, 0, time.Local))
		})
	})

}
