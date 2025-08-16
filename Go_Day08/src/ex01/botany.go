package main

import (
	"fmt"
	"reflect"
	"strings"
)

func describePlant(p interface{}) string {
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)

	if t.Kind() != reflect.Struct {
		return "Input is not a struct"
	}

	var parts []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		name := field.Name

		tagParts := []string{}
		for _, tagKey := range []string{"unit", "color_scheme"} {
			if tagValue := field.Tag.Get(tagKey); tagValue != "" {
				tagParts = append(tagParts, fmt.Sprintf("%s=%s", tagKey, tagValue))
			}
		}
		if len(tagParts) > 0 {
			name += "(" + strings.Join(tagParts, ",") + ")"
		}

		parts = append(parts, fmt.Sprintf("%s:%v", name, value.Interface()))
	}
	return strings.Join(parts, ", ")
}
