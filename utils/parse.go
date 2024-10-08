package utils

import (
	"encoding/json"
	"errors"
	"im/utils/meowlog"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Bind(req *http.Request, obj any) error {
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return BindJson(req, obj)
	}
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return BindForm(req, obj)
	}
	return errors.New("当前方法不支持")
}

func BindJson(req *http.Request, obj any) error {
	s, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(s, obj)
	return err
}

func BindForm(req *http.Request, ptr any) error {
	req.ParseForm()
	logger := meowlog.NewLogger("console", "debug", "")
	logger.Debug("req.Form.Encode():%v\n", req.Form.Encode())
	err := mapForm(ptr, req.Form)
	return err
}
func mapForm(ptr any, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		valField := val.Field(i)
		if !valField.CanSet() {
			continue
		}
		valFieldKind := valField.Kind()
		inputFieldName := typeField.Tag.Get("form")
		if inputFieldName == "" {
			inputFieldName = typeField.Name
			if valFieldKind == reflect.Struct {
				err := mapForm(valField.Addr().Interface(), form)
				if err != nil {
					return err
				}
				continue
			}
		}
		inputValue, ok := form[inputFieldName]
		if !ok {
			continue
		}
		numElems := len(inputValue)
		if valFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := valField.Type().Elem().Kind()
			slice := reflect.MakeSlice(valField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else {
			if _, isTime := valField.Interface().(time.Time); isTime {
				if err := setTimeField(inputValue[0], typeField, valField); err != nil {
					return err
				}
				continue
			}
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], valField); err != nil {
				return err
			}
		}
	}
	return nil
}
func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err != nil {
		return err
	}
	field.SetInt(intVal)
	return nil
}
func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err != nil {
		return err
	}
	field.SetUint(uintVal)
	return nil
}
func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	field.SetBool(boolVal)
	return nil
}
func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err != nil {
		return err
	}
	field.SetFloat(floatVal)
	return nil
}
func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = "2006-01-02 15:04:05"
		val = strings.Replace(val, "/", "-", 0)
		num := len(strings.Split(val, " "))
		if num == 1 {
			val = val + " 00:00:00"
		} else {
			num = len(strings.Split(val, ":"))
			if num == 1 {
				val = val + ":00:00"
			} else if num == 2 {
				val = val + ":00"
			}
		}
	}
	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

// Don't pass in pointers to bind to. Can lead to bugs. See:
// https://github.com/codegangsta/martini-contrib/issues/40
// https://github.com/codegangsta/martini-contrib/pull/34#issuecomment-29683659
func ensureNotPointer(obj interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		panic("Pointers are not accepted as binding models")
	}
}
