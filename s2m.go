package s2m

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const Tag = "s2m"

// Do Switch any value to map[string]any.
// If s is not struct, will format s to string str, and return empty map.
func Do(s any) map[string]any {
	m, _ := DoWithErr(s)
	return m
}

// DoWithErr Switch struct to map[string]any.
// If s is not struct, will return error.
// This switch is not deep.
// If the struct contain other structs, the structs contained will not be witched.
func DoWithErr(s any) (m map[string]any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("s2m: %v", rec)
		}
	}()

	m = map[string]any{}

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
		return DoWithErr(val.Interface())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Anonymous {
			continue
		}

		name := field.Name

		v := val.FieldByName(name)

		if !v.IsValid() {
			continue
		}

		tag, ok := field.Tag.Lookup(Tag)

		if ok {
			if strings.Contains(tag, "omitempty") && v.IsZero() {
				continue
			}
			name = tag
		}

		m[name] = v.Interface()
	}

	return
}
