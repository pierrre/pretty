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
	// ShowKnownTypes shows known types.
	// Default: false.
	ShowKnownTypes bool
	// ShowUnderlyingType shows the underlying type.
	// Default: true.
	ShowUnderlyingType bool
}

// NewTypeValueWriter creates a new [TypeValueWriter] with default values.
func NewTypeValueWriter(vw ValueWriter) *TypeValueWriter {
	return &TypeValueWriter{
		ValueWriter:        vw,
		ShowKnownTypes:     false,
		ShowUnderlyingType: true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !st.KnownType || vw.ShowKnownTypes {
		write.MustString(st.Writer, "[")
		writeType(st.Writer, v.Type())
		write.MustString(st.Writer, "]")
		vw.writeUnderlyingType(st.Writer, v)
		write.MustString(st.Writer, " ")
	}
	knownType := st.KnownType
	st.KnownType = true // The type is known, because we showed it.
	must.Handle(vw.ValueWriter.WriteValue(st, v))
	st.KnownType = knownType
	return true
}

func (vw *TypeValueWriter) writeUnderlyingType(w io.Writer, v reflect.Value) {
	if !vw.ShowUnderlyingType {
		return
	}
	typ := v.Type()
	uTyp := reflectutil.GetUnderlyingType(typ)
	if uTyp != typ {
		write.MustString(w, "(")
		writeType(w, uTyp)
		write.MustString(w, ")")
	}
}

func writeType(w io.Writer, typ reflect.Type) {
	write.MustString(w, reflectutil.TypeFullName(typ))
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
func (vws ByTypeValueWriters) WriteValue(st *State, v reflect.Value) bool {
	if len(vws) == 0 {
		return false
	}
	typ := v.Type()
	vw, ok := vws[typ]
	if !ok {
		return false
	}
	return vw.WriteValue(st, v)
}

// Supports implements [SupportChecker].
func (vws ByTypeValueWriters) Supports(typ reflect.Type) ValueWriter {
	return supportsValueWriter(typ, vws[typ])
}
