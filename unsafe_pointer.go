package pretty

import (
	"reflect"
)

// UnsafePointerValueWriter is a [ValueWriter] that handles unsafe pointer values.
//
// It should be created with [NewUnsafePointerValueWriter].
type UnsafePointerValueWriter struct{}

// NewUnsafePointerValueWriter creates a new [UnsafePointerValueWriter].
func NewUnsafePointerValueWriter() *UnsafePointerValueWriter {
	return &UnsafePointerValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *UnsafePointerValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.UnsafePointer {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	writeUintptr(st.Writer, uintptr(v.UnsafePointer()))
	return true
}
