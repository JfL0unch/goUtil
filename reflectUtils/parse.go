package reflectUtils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
)

// FlatStructFields
// 将结构体的所有字段都返回(会将 Anonymous fields 打平后包含在内)
// 参考: https://stackoverflow.com/questions/24333494/golang-reflection-on-embedded-structs
func FlatStructFields(anonymousField interface{}) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	ifv := reflect.ValueOf(anonymousField)
	ift := reflect.TypeOf(anonymousField)

	for i := 0; i < ift.NumField(); i++ {
		fv := ifv.Field(i)
		ft := ift.Field(i)
		if ft.Type.Kind() == reflect.Struct && ft.Anonymous {
			fields = append(fields, FlatStructFields(fv.Interface())...)
		} else {
			fields = append(fields, ft)
		}

	}

	return fields
}

func isAliasType(zeroVal reflect.Value) bool {
	underlying := zeroVal.Kind()
	instanceType := zeroVal.Type()

	//if underlying.String() == reflect.Struct.String() {
	//	return false
	//}
	return underlying.String() != instanceType.String()
}

func getTimeFromStr(strVal string) (time.Time, error) {
	if strVal == "" || strVal == "0000-00-00 00:00:00" {
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), nil
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", strVal, time.Local)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/1/2 5:4:5", strVal, time.Local)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("01/02/06 15:04", strVal, time.Local)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("2006/01/02 15:04:05", strVal, time.Local)
	if err == nil {
		return t, nil
	}
	t, err = time.ParseInLocation("1/2/06 15:4", strVal, time.Local)
	if err == nil {
		return t, nil
	}
	// 2023-07-04T09:36:33
	t, err = time.ParseInLocation("2006-01-02T15:04:05", strVal, time.Local)
	if err == nil {

		return t, nil
	}
	// 2023-07-04T09:36:33.2961605775
	t, err = time.ParseInLocation("2006-01-02T15:14:15.0000000000", strVal, time.Local)
	if err == nil {
		return t, nil
	}

	return time.Date(1971, 1, 1, 0, 0, 0, 0, time.Local), nil

}
func getInstanceOfAliasType(aliasTypeZeroVal reflect.Value, strVal string) (reflect.Value, error) {

	underlyingKind := aliasTypeZeroVal.Kind()
	switch underlyingKind {
	case reflect.Struct:
		if !json.Valid([]byte(strVal)) { // 无效json
			switch aliasTypeZeroVal.Type().Name() {
			case "Time":
				t, err := getTimeFromStr(strVal)
				if err != nil {
					return reflect.ValueOf(t), errors.Wrapf(err, "ParseInLocation(%s)", strVal)
				}
				return reflect.ValueOf(t), nil
			default:
				return reflect.ValueOf(""), fmt.Errorf("strVal(%s) not json-format", strVal)
			}
		} else { // 有效json
			v := reflect.New(aliasTypeZeroVal.Type())
			if strVal != "" {
				err := json.Unmarshal([]byte(strVal), v.Interface())
				if err != nil {
					return reflect.Value{}, errors.Wrapf(err, "json val=%s", strVal)
				}
			}

			return v.Elem(), nil
		}
	default:
		v, err := getInstance(aliasTypeZeroVal, strVal)
		if err != nil {
			return aliasTypeZeroVal, errors.Wrapf(err, "getInstance(%v,%s)", aliasTypeZeroVal, strVal)
		}
		return v.Convert(aliasTypeZeroVal.Type()), nil
	}

}

func getInstance(instanceZeroVal reflect.Value, strVal string) (reflect.Value, error) {
	underlyingKind := instanceZeroVal.Kind()

	//var x reflect.Value
	switch underlyingKind {
	case reflect.String:
		return reflect.ValueOf(strVal), nil
	case reflect.Int:
		if strVal == "" {
			return reflect.ValueOf(0), nil
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(0), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(int(i)), nil
	case reflect.Int8:
		if strVal == "" {
			return reflect.ValueOf(int8(0)), nil
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(int8(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(int8(i)), nil
	case reflect.Int16:
		if strVal == "" {
			return reflect.ValueOf(int16(0)), nil
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(int16(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(int16(i)), nil
	case reflect.Int32:
		if strVal == "" {
			return reflect.ValueOf(int32(0)), nil
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(int32(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(int32(i)), nil
	case reflect.Int64:
		if strVal == "" {
			return reflect.ValueOf(int64(0)), nil
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(int64(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(i), nil
	case reflect.Uint:
		if strVal == "" {
			return reflect.ValueOf(uint(0)), nil
		}
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(uint(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(uint(i)), nil
	case reflect.Uint8:
		if strVal == "" {
			return reflect.ValueOf(uint8(0)), nil
		}
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(uint8(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(uint8(i)), nil
	case reflect.Uint16:
		if strVal == "" {
			return reflect.ValueOf(uint16(0)), nil
		}
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(uint16(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(uint16(i)), nil
	case reflect.Uint32:
		if strVal == "" {
			return reflect.ValueOf(uint32(0)), nil
		}
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(uint32(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(uint32(i)), nil
	case reflect.Uint64:
		if strVal == "" {
			return reflect.ValueOf(uint64(0)), nil
		}
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return reflect.ValueOf(uint64(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(i), nil
	case reflect.Uintptr:
		if strVal == "" {
			dummy := uint(0)
			return reflect.ValueOf(&dummy), nil
		}
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			dummy := uint(0)
			return reflect.ValueOf(&dummy), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(&i), nil
	case reflect.Float32:
		if strVal == "" {
			return reflect.ValueOf(float32(0)), nil
		}
		i, err := strconv.ParseFloat(strVal, 64)
		if err == nil {
			return reflect.ValueOf(float32(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(float32(i)), nil

	case reflect.Float64:
		if strVal == "" {
			return reflect.ValueOf(float64(0)), nil
		}
		i, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return reflect.ValueOf(float64(0)), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(i), nil
	case reflect.Complex64:
		if strVal == "" {
			return reflect.ValueOf(complex64(0)), nil
		}
		i, err := strconv.ParseComplex(strVal, 64)
		if err != nil {
			return reflect.ValueOf(strVal), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(complex64(i)), nil
	case reflect.Complex128:
		if strVal == "" {
			return reflect.ValueOf(complex128(0)), nil
		}
		i, err := strconv.ParseComplex(strVal, 64)
		if err != nil {
			return reflect.ValueOf(strVal), errors.Wrapf(err, "val=%s", strVal)
		}
		return reflect.ValueOf(i), nil
	case reflect.Bool:
		if strVal == "true" {
			return reflect.ValueOf(true), nil
		} else {
			return reflect.ValueOf(false), nil
		}
	case reflect.Struct:
		if strVal == "" {
			v := reflect.New(instanceZeroVal.Type())
			return v.Elem(), nil
		}

		if false { // isAliasType(instanceZeroVal)
			return getInstanceOfAliasType(instanceZeroVal, strVal)
		} else {
			v := reflect.New(instanceZeroVal.Type())
			if strVal != "" {
				err := json.Unmarshal([]byte(strVal), v.Interface())
				if err != nil {
					return reflect.Value{}, errors.Wrapf(err, "json val=%s", strVal)
				}
			}

			return v.Elem(), nil
		}

	case reflect.Map, reflect.Slice, reflect.Array: // strVal should be json
		v := reflect.New(instanceZeroVal.Type())

		if strVal == "" {
			return v.Elem(), nil
		}
		if json.Valid([]byte(strVal)) {
			err := json.Unmarshal([]byte(strVal), v.Interface())
			if err != nil {
				return reflect.Value{}, errors.Wrapf(err, "val=%s", strVal)
			}
		}

		return v.Elem(), nil

	case reflect.Pointer:
		if strVal == "" || strVal == "null" {
			return instanceZeroVal, nil
		}

		instanceVal := instanceZeroVal.Elem()
		if !instanceVal.IsValid() {
			return instanceZeroVal, nil
		}

		instanceType := instanceVal.Type()

		gotVal, err := ParseStrToInstance(instanceVal, strVal)

		if err == nil {
			fieldPtr := reflect.New(instanceType)

			field := fieldPtr.Elem() // ！！！ important: only Elem() CanSet !!!

			field.Set(gotVal)

			return fieldPtr, nil
		} else {
			return reflect.Value{}, errors.Wrapf(err, "ParseStrToInstance(%v,%s)", instanceVal, strVal)
		}
	default:
		return reflect.Value{}, errors.Errorf("ParseStrToInstance(%v,%s) not support %s", instanceZeroVal, strVal, underlyingKind.String())
	}
}

// ParseStrToInstance
// 把字符串翻译为指定类型
//
// supported:
// struct
// standard type(int,string,bool...)
//
// not supported:
// reflect.Chan
// reflect.Func
// reflect.Interface
// reflect.UnsafePointer
func ParseStrToInstance(zeroVal reflect.Value, strVal string) (reflect.Value, error) {

	var instance reflect.Value
	var err error

	if isAliasType(zeroVal) {
		return getInstanceOfAliasType(zeroVal, strVal)
	} else {
		instance, err = getInstance(zeroVal, strVal)
		if err != nil {
			return instance, err
		}
		return instance, err
	}
}
