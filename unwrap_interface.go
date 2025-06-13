package pretty

import (
	"reflect"
)

// UnwrapInterfaceWriter is a [ValueWriter] that unwraps interface values.
//
// It should be created with [NewUnwrapInterfaceWriter].
type UnwrapInterfaceWriter struct {
	ValueWriter
}

// NewUnwrapInterfaceWriter creates a new [UnwrapInterfaceWriter].
func NewUnwrapInterfaceWriter(vw ValueWriter) *UnwrapInterfaceWriter {
	return &UnwrapInterfaceWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *UnwrapInterfaceWriter) WriteValue(st *State, v reflect.Value) bool {
	v, isNil := vw.unwrapInterface(st, v)
	return isNil || vw.ValueWriter.WriteValue(st, v)
}

func (vw *UnwrapInterfaceWriter) unwrapInterface(st *State, v reflect.Value) (_ reflect.Value, isNil bool) {
	if v.Kind() == reflect.Interface {
		if checkNil(st.Writer, v) {
			return reflect.Value{}, true
		}
		v = v.Elem()
		st.KnownType = false // We want to show the type of the value.
	}
	return v, false
}
