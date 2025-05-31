// Package itfassert provides utilities to assert that a reflect.Value implements a specific interface type.
package itfassert

import (
	"reflect"
)

// Assert checks if the given reflect.Value can be asserted to type T.
// It returns the value of type T and a boolean indicating whether the assertion was successful.
func Assert[T any](v reflect.Value) (T, bool) {
	var zero T
	if !v.CanInterface() {
		return zero, false
	}
	kind := v.Kind()
	if kind == reflect.Pointer && v.IsNil() {
		// Prevents calling methods on nil pointers.
		return zero, false
	}
	vi, ok := v.Interface().(T)
	return vi, ok
}
