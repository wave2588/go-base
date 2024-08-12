package utils

import "reflect"

func RecursiveIndirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func RecursiveIndirectType(p reflect.Type) reflect.Type {
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	return p
}
