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
	v, isNil := unwrapInterface(st, v)
	return isNil || vw.ValueWriter.WriteValue(st, v)
}

func unwrapInterface(st *State, v reflect.Value) (_ reflect.Value, isNil bool) {
	if v.Kind() == reflect.Interface {
		if checkNil(st.Writer, v) {
			return reflect.Value{}, true
		}
		v = v.Elem()
		st.KnownType = false // We want to show the type of the value.
	}
	return v, false
}
