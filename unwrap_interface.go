package pretty

import (
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
func (vw *UnwrapInterfaceValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() == reflect.Interface {
		if checkNil(st.Writer, v) {
			return true
		}
		v = v.Elem()
		defer st.setRestoreKnownType(false)()
	}
	return vw.ValueWriter(st, v)
}
