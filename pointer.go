package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/must"
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
func (vw *PointerValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	infos{
		showAddr: vw.ShowAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(st)
	writeArrow(st.Writer)
	must.Handle(vw.ValueWriter.WriteValue(st, v.Elem()))
	return true
}
