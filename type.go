package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// TypeValueWriter is a [ValueWriter] that writes the type of the value.
//
// It should be created with [NewTypeValueWriter].
type TypeValueWriter struct {
	ValueWriter
	// Stringer converts the [reflect.Type] to a string.
	// Default: [reflectutil.TypeFullName].
	Stringer func(reflect.Type) string
	// ShowKnownTypes shows known types.
	// Default: false.
	ShowKnownTypes bool
	// ShowBaseType shows the base type.
	// Default: true.
	ShowBaseType bool
}

// NewTypeValueWriter creates a new [TypeValueWriter] with default values.
func NewTypeValueWriter(vw ValueWriter) *TypeValueWriter {
	return &TypeValueWriter{
		ValueWriter:    vw,
		Stringer:       reflectutil.TypeFullName,
		ShowKnownTypes: false,
		ShowBaseType:   true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !st.KnownType || vw.ShowKnownTypes {
		write.MustString(st.Writer, "[")
		vw.writeType(st.Writer, v.Type())
		write.MustString(st.Writer, "]")
		vw.writeBaseType(st.Writer, v)
		write.MustString(st.Writer, " ")
	}
	knownType := st.KnownType
	st.KnownType = true // The type is known, because we showed it.
	must.Handle(vw.ValueWriter.WriteValue(st, v))
	st.KnownType = knownType
	return true
}

func (vw *TypeValueWriter) writeType(w io.Writer, typ reflect.Type) {
	write.MustString(w, vw.Stringer(typ))
}

func (vw *TypeValueWriter) writeBaseType(w io.Writer, v reflect.Value) {
	if !vw.ShowBaseType {
		return
	}
	typ := v.Type()
	baseType := reflectutil.GetBaseType(typ)
	if baseType == nil || baseType == typ {
		return
	}
	write.MustString(w, "(")
	vw.writeType(w, baseType)
	write.MustString(w, ")")
}

// ByTypeValueWriters is a [ValueWriter] that selects a [ValueWriter] by [reflect.Type].
//
// It should be created with [NewByTypeValueWriters].
type ByTypeValueWriters map[reflect.Type]ValueWriter

// NewByTypeValueWriters creates a new [ByTypeValueWriters].
func NewByTypeValueWriters() ByTypeValueWriters {
	return make(ByTypeValueWriters)
}

// WriteValue implements [ValueWriter].
func (vw ByTypeValueWriters) WriteValue(st *State, v reflect.Value) bool {
	if len(vw) == 0 {
		return false
	}
	typ := v.Type()
	w, ok := vw[typ]
	if !ok {
		return false
	}
	return w.WriteValue(st, v)
}
