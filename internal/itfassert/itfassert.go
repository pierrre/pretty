// Package itfassert provides utilities to assert that a reflect.Value implements a specific interface type.
package itfassert

import (
	"reflect"
)

// Assert checks if the given [reflect.Value] can be asserted to type T.
// It returns the value of type T and a boolean indicating whether the assertion was successful.
// If the value is a pointer and is nil, it returns the zero value of T and false to prevent calling methods on nil pointers.
func Assert[T any](v reflect.Value) (T, bool) {
	// TODO: use reflect.TypeAssert when available in Go 1.25.
	var zero T
	if !v.CanInterface() {
		return zero, false
	}
	kind := v.Kind()
	if kind == reflect.Pointer && v.IsNil() {
		// Prevents calling methods on nil pointers.
		return zero, false
	}
	return reflect.TypeAssert[T](v)
}
