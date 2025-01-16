package pretty

import (
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
func (vw *InterfaceValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Interface {
		return false
	}
	writeArrow(st.Writer)
	if checkNil(st.Writer, v) {
		return true
	}
	defer st.SetRestoreKnownType(false)() // We want to show the type of the value.
	mustHandle(vw.ValueWriter(st, v.Elem()))
	return true
}
