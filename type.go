package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// TypeWriter is a [ValueWriter] that writes the type of the value.
//
// It should be created with [NewTypeWriter].
type TypeWriter struct {
	ValueWriter
	// ShowKnownTypes shows known types.
	// Default: false.
	ShowKnownTypes bool
	// ShowUnderlyingType shows the underlying type.
	// Default: true.
	ShowUnderlyingType bool
}

// NewTypeWriter creates a new [TypeWriter] with default values.
func NewTypeWriter(vw ValueWriter) *TypeWriter {
	return &TypeWriter{
		ValueWriter:        vw,
		ShowKnownTypes:     false,
		ShowUnderlyingType: true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeWriter) WriteValue(st *State, v reflect.Value) bool {
	knownType := vw.writeType(st, v)
	vw.ValueWriter.WriteValue(st, v)
	vw.postType(st, knownType)
	return true
}

func (vw *TypeWriter) writeType(st *State, v reflect.Value) (knownType bool) {
	if !st.KnownType || vw.ShowKnownTypes {
		typ := v.Type()
		st.Writer.AppendByte('[')
		st.Writer.AppendString(reflectutil.TypeFullName(typ))
		st.Writer.AppendByte(']')
		if vw.ShowUnderlyingType {
			uTyp := reflectutil.GetUnderlyingType(typ)
			if uTyp != typ {
				st.Writer.AppendByte('(')
				st.Writer.AppendString(reflectutil.TypeFullName(uTyp))
				st.Writer.AppendByte(')')
			}
		}
		st.Writer.AppendByte(' ')
	}
	knownType = st.KnownType
	st.KnownType = true // The type is known, because we showed it.
	return knownType
}

func (vw *TypeWriter) postType(st *State, knownType bool) {
	st.KnownType = knownType
}

// ByTypeWriters is a [ValueWriter] that selects a [ValueWriter] by [reflect.Type].
//
// It should be created with [NewByTypeWriters].
type ByTypeWriters map[reflect.Type]ValueWriter

// NewByTypeWriters creates a new [ByTypeWriters].
func NewByTypeWriters() ByTypeWriters {
	return make(ByTypeWriters)
}

// WriteValue implements [ValueWriter].
func (vws ByTypeWriters) WriteValue(st *State, v reflect.Value) bool {
	if len(vws) != 0 {
		typ := v.Type()
		vw, ok := vws[typ]
		if ok {
			return vw.WriteValue(st, v)
		}
	}
	return false
}

// Supports implements [SupportChecker].
func (vws ByTypeWriters) Supports(typ reflect.Type) ValueWriter {
	return supportsValueWriter(typ, vws[typ])
}
