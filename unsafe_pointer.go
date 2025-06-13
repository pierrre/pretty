package pretty

import (
	"reflect"
)

// UnsafePointerWriter is a [ValueWriter] that handles unsafe pointer values.
//
// It should be created with [NewUnsafePointerWriter].
type UnsafePointerWriter struct{}

// NewUnsafePointerWriter creates a new [UnsafePointerWriter].
func NewUnsafePointerWriter() *UnsafePointerWriter {
	return &UnsafePointerWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *UnsafePointerWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.UnsafePointer {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	writeUintptr(st.Writer, uintptr(v.UnsafePointer()))
	return true
}

// Supports implements [SupportChecker].
func (vw *UnsafePointerWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.UnsafePointer {
		res = vw
	}
	return res
}
