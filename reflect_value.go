package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
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
	if v.Kind() != reflect.Struct || v.Type() != reflectValueType {
		return false
	}
	if !v.CanInterface() {
		write.MustString(st.Writer, "<unexported>")
		return true
	}
	rv := v.Interface().(reflect.Value) //nolint:forcetypeassert // Checked above.
	writeArrow(st.Writer)
	if checkInvalidNil(st.Writer, rv) {
		return true
	}
	st.KnownType = false // We want to show the type of the value.
	must.Handle(vw.ValueWriter.WriteValue(st, rv))
	return true
}

// Supports implements [SupportChecker].
func (vw *ReflectValueWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == reflectValueType {
		res = vw
	}
	return res
}
