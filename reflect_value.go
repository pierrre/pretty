package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal"
	"github.com/pierrre/pretty/internal/must"
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
func (vw *ReflectValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Type() != reflectValueType {
		return false
	}
	if !v.CanInterface() {
		internal.MustWriteString(st.Writer, "<unexported>")
		return true
	}
	rv := v.Interface().(reflect.Value) //nolint:forcetypeassert // Checked above.
	writeArrow(st.Writer)
	if checkInvalid(st.Writer, rv) {
		return true
	}
	defer st.SetRestoreKnownType(false)() // We want to show the type of the value.
	must.Handle(vw.ValueWriter(st, rv))
	return true
}
