package s2m

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const Tag = "s2m"

// To Switch any value to map[string]any.
// If s is not struct, will format s to string str, and return empty map.
func To(s any) map[string]any {
	m, _ := ToWithErr(s)
	return m
}

// ToWithErr Switch struct to map[string]any.
// If s is not struct, will return error.
// This switch is not deep.
// If the struct contain other structs, the structs contained will not be witched.
func ToWithErr(s any) (m map[string]any, err error) {
	return doWithErr[any](s, false)
}

func ToStrMap(s any) map[string]string {
	m, _ := ToStrMapWithErr(s)
	return m
}

func ToStrMapWithErr(s any) (map[string]string, error) {
	return doWithErr[string](s, true)
}

func doWithErr[V any](s any, isValStr bool) (m map[string]V, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("s2m: %v", rec)
		}
	}()
	m = map[string]V{}
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)
	kind := val.Kind()
	if kind != reflect.Struct {
		err = fmt.Errorf("s2m: input kind %v is not a struct", kind)
		return
	}
	if !val.IsValid() {
		err = errors.New("s2m: input is invalid")
		return
	}
	if val.IsNil() {
		err = errors.New("s2m: input is nil")
		return
	}
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.UnsafePointer {
		return doWithErr[V](val.Interface(), isValStr)
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			continue
		}
		name := field.Name
		v := val.FieldByName(name)
		if !v.IsValid() {
			err = fmt.Errorf("s2m: field %v is invalid", name)
			return
		}
		tag, ok := field.Tag.Lookup(Tag)
		if ok {
			if strings.Contains(tag, "omitempty") && v.IsZero() {
				continue
			}
			name = tag
		}
		if isValStr {
			var str string
			str, err = formatAtom(v)
			if err != nil {
				return
			}
			m[name] = (any)(str).(V)
		} else {
			m[name] = v.Interface()
		}
	}
	return
}

func formatAtom(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float64, reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.String:
		return strconv.Quote(v.String()), nil
	case reflect.Slice, reflect.Map, reflect.Struct, reflect.Array:
		data, err := json.Marshal(v.Interface())
		return string(data), err
	case reflect.Chan, reflect.Func:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16), nil
	case reflect.Ptr:
		return formatAtom(v.Elem())
	case reflect.Invalid:
		return "invalid", nil
	default: // reflect.Interface
		return v.Type().String() + " value", errors.New("s2m: interface")
	}
}
