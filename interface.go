package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/must"
)

// InterfaceWriter is a [ValueWriter] that handles interface values.
//
// It should be created with [NewInterfaceWriter].
type InterfaceWriter struct {
	ValueWriter
}

// NewInterfaceWriter creates a new [InterfaceWriter].
func NewInterfaceWriter(vw ValueWriter) *InterfaceWriter {
	return &InterfaceWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *InterfaceWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Interface {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	writeArrow(st.Writer)
	st.KnownType = false // We want to show the type of the value.
	must.Handle(vw.ValueWriter.WriteValue(st, v.Elem()))
	return true
}

// Supports implements [SupportChecker].
func (vw *InterfaceWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Interface {
		res = vw
	}
	return res
}
