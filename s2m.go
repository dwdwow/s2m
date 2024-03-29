package s2m

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const Tag = "s2m"

// ToMapNoErr Switch any value toMap map[string]any.
// If s is not struct, will format s toMap string str, and return empty map.
func ToMapNoErr(s any) map[string]any {
	m, _ := ToMap(s)
	return m
}

// ToMap Switch struct toMap map[string]any.
// If s is not struct, will return error.
// This switch is not deep.
// If the struct contain other structs, the structs contained will not be witched.
func ToMap(s any) (m map[string]any, err error) {
	return toMap[any](s, false)
}

func ToStrMapNoErr(s any) map[string]string {
	m, _ := ToStrMap(s)
	return m
}

func ToStrMap(s any) (map[string]string, error) {
	return toMap[string](s, true)
}

func toMap[V any](s any, isValStr bool) (m map[string]V, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("s2m: %v", rec)
		}
	}()
	m = map[string]V{}
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)
	kind := val.Kind()
	// s may be nil, invalid value
	if !val.IsValid() {
		return m, nil
	}
	if kind != reflect.Struct {
		err = fmt.Errorf("s2m: input kind %v is not a struct", kind)
		return
	}
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.UnsafePointer {
		return toMap[V](val.Interface(), isValStr)
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Name
		v := val.FieldByName(name)
		if !v.IsValid() {
			err = fmt.Errorf("s2m: field %v is invalid", name)
			m = map[string]V{}
			return
		}
		tag, ok := field.Tag.Lookup(Tag)
		if ok {
			// field tag name may be omitempty, so can not use strings.Contains(tag, "omitempty").
			tags := strings.Split(tag, ",")
			if len(tags) > 1 &&
				strings.Contains(tags[1], "omitempty") &&
				v.IsZero() {
				continue
			}
			name = strings.Trim(tags[0], " ")
		}
		if isValStr {
			var str string
			str, err = formatAtom(v)
			if err != nil {
				return
			}
			m[name] = (any)(str).(V)
		} else {
			m[name] = v.Interface().(V)
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
		return v.String(), nil
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
		return v.Type().String() + " value", nil
	}
}
