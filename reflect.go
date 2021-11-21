package ddbstruct

import (
	"reflect"
)

func maybealloc(v reflect.Value) reflect.Value {
	if v.Type().Kind() != reflect.Pointer {
		return v.Addr()
	}
	if v.IsNil() {
		n := reflect.New(v.Type().Elem())
		v.Set(n)
		return n
	}
	return v
}

func getF(d interface{}, f int) reflect.Value {
	return maybealloc(reflect.ValueOf(d).Elem().Field(f)).Elem()
}
