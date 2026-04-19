package pretty

import (
	"reflect"
)

// PointerWriter is a [ValueWriter] that handles pointer values.
//
// It should be created with [NewPointerWriter].
type PointerWriter struct {
	ValueWriter
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
}

// NewPointerWriter creates a new [PointerWriter] with default values.
func NewPointerWriter(vw ValueWriter) *PointerWriter {
	return &PointerWriter{
		ValueWriter: vw,
		ShowAddr:    true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *PointerWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	if checkNil(st, v) {
		return true
	}
	infos{
		showAddr: vw.ShowAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(st)
	writeArrow(st)
	vw.ValueWriter.WriteValue(st, v.Elem())
	return true
}

// Supports implements [SupportChecker].
func (vw *PointerWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Pointer {
		res = vw
	}
	return res
}
