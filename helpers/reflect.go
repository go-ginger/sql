package helpers

import "reflect"

type Value struct {
	reflect.Value
}

func CreateArray(t reflect.Type) reflect.Value {
	var arrayType reflect.Type
	arrayType = reflect.SliceOf(t)
	return reflect.Zero(arrayType)
}
