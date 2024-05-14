package ui

import (
	"fmt"
	"reflect"
	"strings"
)

func ToTable[T any](s T) string {
	var result strings.Builder

	v := reflect.ValueOf(s)

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			name := t.Field(i).Name
			writeValue(&result, name, field)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			writeValue(&result, fmt.Sprintf("%v", key.Interface()), value)
		}
	default:
		result.WriteString(fmt.Sprintf("Unsupported type: %s\n", v.Kind()))
	}

	return result.String()
}

func writeValue(result *strings.Builder, name string, value reflect.Value) {
	switch value.Kind() {
	case reflect.String:
		if value.String() != "" {
			result.WriteString(fmt.Sprintf("%s: %s\n", name, value.String()))
		}
	case reflect.Int:
		if value.Int() != 0 {
			result.WriteString(fmt.Sprintf("%s: %d\n", name, value.Int()))
		}
	case reflect.Map:
		if value.Len() > 0 {
			result.WriteString(fmt.Sprintf("%s:\n", name))
			for _, key := range value.MapKeys() {
				writeValue(result, fmt.Sprintf("  %s", key), value.MapIndex(key))
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			writeValue(result, fmt.Sprintf("%s[%d]", name, i), value.Index(i))
		}
	case reflect.Struct:
		t := value.Type()
		for i := 0; i < value.NumField(); i++ {
			writeValue(result, fmt.Sprintf("%s.%s", name, t.Field(i).Name), value.Field(i))
		}
	}
}
