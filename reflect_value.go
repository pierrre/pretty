package pretty

import (
	"io"
	"reflect"
)

var reflectValueType = reflect.TypeFor[reflect.Value]()

// ReflectValueWriter is a [ValueWriter] that handles [reflect.Value].
//
// It should be created with [NewReflectValueWriter].
type ReflectValueWriter struct {
	ValueWriter
}

// NewReflectValueWriter creates a new [ReflectValueWriter].
func NewReflectValueWriter(vw ValueWriter) *ReflectValueWriter {
	return &ReflectValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ReflectValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Type() != reflectValueType {
		return false
	}
	if !v.CanInterface() {
		writeString(w, "<unexported>")
		return true
	}
	rv := v.Interface().(reflect.Value) //nolint:forcetypeassert // Checked above.
	writeArrow(w)
	if checkInvalid(w, rv) {
		return true
	}
	st.KnownType = false
	mustHandle(vw.ValueWriter(w, st, rv))
	return true
}
