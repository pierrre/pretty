package pretty

import (
	"io"
	"reflect"
)

// PointerValueWriter is a [ValueWriter] that handles pointer values.
//
// It should be created with [NewPointerValueWriter].
type PointerValueWriter struct {
	ValueWriter
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
}

// NewPointerValueWriter creates a new [PointerValueWriter] with default values.
func NewPointerValueWriter(vw ValueWriter) *PointerValueWriter {
	return &PointerValueWriter{
		ValueWriter: vw,
		ShowAddr:    true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *PointerValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	infos{
		showAddr: vw.ShowAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(w)
	writeArrow(w)
	mustHandle(vw.ValueWriter(w, st, v.Elem()))
	return true
}
