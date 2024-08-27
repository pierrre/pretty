package pretty

import (
	"io"
	"reflect"
)

// UnwrapInterfaceValueWriter is a [ValueWriter] that unwraps interface values.
//
// It should be created with [NewUnwrapInterfaceValueWriter].
type UnwrapInterfaceValueWriter struct {
	ValueWriter
}

// NewUnwrapInterfaceValueWriter creates a new [UnwrapInterfaceValueWriter].
func NewUnwrapInterfaceValueWriter(vw ValueWriter) *UnwrapInterfaceValueWriter {
	return &UnwrapInterfaceValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *UnwrapInterfaceValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() == reflect.Interface {
		if checkNil(w, v) {
			return true
		}
		v = v.Elem()
		st.KnownType = false
	}
	return vw.ValueWriter(w, st, v)
}
