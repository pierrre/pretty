package pretty

import (
	"io"
	"reflect"
)

// InterfaceValueWriter is a [ValueWriter] that handles interface values.
//
// It should be created with [NewInterfaceValueWriter].
type InterfaceValueWriter struct {
	ValueWriter
}

// NewInterfaceValueWriter creates a new [InterfaceValueWriter].
func NewInterfaceValueWriter(vw ValueWriter) *InterfaceValueWriter {
	return &InterfaceValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *InterfaceValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() != reflect.Interface {
		return false
	}
	writeArrow(w)
	if checkNil(w, v) {
		return true
	}
	st.KnownType = false
	mustHandle(vw.ValueWriter(w, st, v.Elem()))
	return true
}
