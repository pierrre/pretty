package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/syncutil"
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
	// ShowBaseType shows the base type.
	// Default: true.
	ShowBaseType bool
}

// NewTypeValueWriter creates a new [TypeValueWriter] with default values.
func NewTypeValueWriter(vw ValueWriter) *TypeValueWriter {
	return &TypeValueWriter{
		ValueWriter:    vw,
		ShowKnownTypes: false,
		ShowBaseType:   true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !st.KnownType || vw.ShowKnownTypes {
		write.MustString(st.Writer, "[")
		writeType(st.Writer, v.Type())
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

func (vw *TypeValueWriter) writeBaseType(w io.Writer, v reflect.Value) {
	if !vw.ShowBaseType {
		return
	}
	typ := v.Type()
	baseType := reflectutil.GetBaseType(typ)
	if baseType != typ {
		write.MustString(w, "(")
		writeType(w, baseType)
		write.MustString(w, ")")
	}
}

func writeType(w io.Writer, typ reflect.Type) {
	write.MustString(w, reflectutil.TypeFullName(typ))
}

type ByTypeValueWriter struct {
	cache           syncutil.Map[reflect.Type, ValueWriterFunc]
	ValueWriters    map[reflect.Type]ValueWriter
	SupportCheckers []TypeSupportChecker
}

func NewByTypeValueWriter() *ByTypeValueWriter {
	return &ByTypeValueWriter{
		ValueWriters: make(map[reflect.Type]ValueWriter),
	}
}

// WriteValue implements [ValueWriter].
func (vw *ByTypeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	typ := v.Type()
	f, ok := vw.cache.Load(typ)
	if !ok {
		f = vw.getValueWriterFunc(typ)
		vw.cache.Store(typ, f)
	}
	if f == nil {
		return false
	}
	return f(st, v)
}

func (vw *ByTypeValueWriter) getValueWriterFunc(typ reflect.Type) ValueWriterFunc {
	w, ok := vw.ValueWriters[typ]
	if ok {
		return w.WriteValue
	}
	for _, sc := range vw.SupportCheckers {
		f := sc.SupportsType(typ)
		if f != nil {
			return f
		}
	}
	return nil
}

type TypeSupportChecker interface {
	SupportsType(typ reflect.Type) ValueWriterFunc
}
