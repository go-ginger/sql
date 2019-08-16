package helpers

import "reflect"

type Value struct {
	reflect.Value
}

func CreateArray(t reflect.Type, length int) reflect.Value {
	var arrayType reflect.Type
	arrayType = reflect.ArrayOf(length, t)
	return reflect.Zero(arrayType)
}
