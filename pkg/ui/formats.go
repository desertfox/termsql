package ui

import (
	"fmt"
	"reflect"
	"strings"
)

func ToTwoLineString[T any](s T) string {
	var (
		columns []string
		values  []string
	)

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)

			columns = append(columns, t.Field(i).Name)
			values = append(values, strings.ReplaceAll(fmt.Sprintf("%v", field.Interface()), "\n", ""))
		}
	}

	maxLengths := make([]int, len(columns))
	for i := 0; i < len(columns); i++ {
		maxLengths[i] = max(len(columns[i]), len(values[i]))
	}

	for i := 0; i < len(columns); i++ {
		columns[i] = padRight(columns[i], " ", maxLengths[i])
		values[i] = padRight(values[i], " ", maxLengths[i])
	}

	return fmt.Sprintf("%s\n%s\n", strings.Join(columns, " | "), strings.Join(values, " | "))
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
