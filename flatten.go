package main

import (
	"reflect"
)

func flatten(v ...interface{}) []string {
	args := flattenDeep(nil, reflect.ValueOf(v))
	var strings []string

	for _, i := range args {
		strings = append(strings, i.(string))
	}
	return strings
}

func flattenDeep(args []interface{}, v reflect.Value) []interface{} {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			args = flattenDeep(args, v.Index(i))
		}
	} else {
		args = append(args, v.Interface())
	}
	return args
}
