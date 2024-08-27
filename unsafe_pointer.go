package pretty

import (
	"io"
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
func (vw *UnsafePointerValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() != reflect.UnsafePointer {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	writeUintptr(w, uintptr(v.UnsafePointer()))
	return true
}
