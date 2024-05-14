package ui

import (
	"fmt"
	"reflect"
	"strings"
)

func ToTable[T any](s T) string {
	var result strings.Builder

	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Name

		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				result.WriteString(fmt.Sprintf("%s: %s\n", name, field.String()))
			}
		case reflect.Int:
			if field.Int() != 0 {
				result.WriteString(fmt.Sprintf("%s: %d\n", name, field.Int()))
			}
		}
	}

	return result.String()
}
