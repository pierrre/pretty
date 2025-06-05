// Package protobuf provides a [pretty.ValueWriter] for protobuf messages.
package protobuf

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var messageImplementsCache = reflectutil.NewImplementsCacheFor[protoreflect.ProtoMessage]()

// ConfigureDefault configures [pretty.DefaultCommonValueWriter] with [ConfigureCommonValueWriter].
func ConfigureDefault() {
	ConfigureCommonValueWriter(pretty.DefaultCommonValueWriter)
}

// ConfigureCommonValueWriter configures a [pretty.CommonValueWriter] with a [ValueWriter].
func ConfigureCommonValueWriter(vw *pretty.CommonValueWriter) {
	vw.ValueWriters = append(vw.ValueWriters, NewValueWriter(vw))
}

// ValueWriter is a [pretty.ValueWriter] that handles protobuf messages.
//
// It should be created with [NewValueWriter].
type ValueWriter struct {
	pretty.ValueWriter
}

// NewValueWriter creates a new [ValueWriter].
func NewValueWriter(vw pretty.ValueWriter) *ValueWriter {
	return &ValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [pretty.ValueWriter].
func (vw *ValueWriter) WriteValue(st *pretty.State, v reflect.Value) bool {
	if !messageImplementsCache.ImplementedBy(v.Type()) {
		return false
	}
	pm, ok := itfassert.Assert[protoreflect.ProtoMessage](v)
	if !ok {
		return false
	}
	m := pm.ProtoReflect()
	vw.writeMessage(st, m)
	return true
}

func (vw *ValueWriter) writeMessage(st *pretty.State, m protoreflect.Message) {
	write.MustString(st.Writer, "{")
	fs := m.Descriptor().Fields()
	l := fs.Len()
	hasFields := false
	st.IndentLevel++
	for i := range l {
		fd := fs.Get(i)
		if fd.ContainingOneof() != nil && !m.Has(fd) {
			continue
		}
		if !hasFields {
			write.MustString(st.Writer, "\n")
			hasFields = true
		}
		st.WriteIndent()
		write.MustString(st.Writer, string(fd.Name()))
		write.MustString(st.Writer, ": ")
		st.KnownType = false // We want to show the types of fields and values.
		must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(vw.getInterface(m.Get(fd), fd))))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	if hasFields {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
}

func (vw *ValueWriter) getInterface(v protoreflect.Value, fd protoreflect.FieldDescriptor) any {
	itf := v.Interface()
	switch itf := itf.(type) {
	case protoreflect.Message:
		return itf.Interface()
	case protoreflect.List:
		return vw.getList(itf, fd)
	case protoreflect.Map:
		return vw.getMap(itf, fd)
	case protoreflect.EnumNumber:
		return vw.getEnum(itf, fd)
	}
	return itf
}

func (vw *ValueWriter) getList(l protoreflect.List, fd protoreflect.FieldDescriptor) any {
	// TODO create typed slice
	res := make([]any, l.Len())
	for i := range l.Len() {
		res[i] = vw.getInterface(l.Get(i), fd)
	}
	return res
}

func (vw *ValueWriter) getMap(m protoreflect.Map, fd protoreflect.FieldDescriptor) any {
	// TODO create typed map
	res := make(map[any]any, m.Len())
	for key, value := range m.Range {
		res[vw.getInterface(key.Value(), fd.MapKey())] = vw.getInterface(value, fd.MapValue())
	}
	return res
}

func (vw *ValueWriter) getEnum(e protoreflect.EnumNumber, fd protoreflect.FieldDescriptor) EnumValue {
	res := EnumValue{
		Number: int32(e),
	}
	ed := fd.Enum().Values().ByNumber(e)
	if ed != nil {
		res.Name = string(ed.Name())
	}
	return res
}

// Supports implements [pretty.SupportChecker].
func (vw *ValueWriter) Supports(typ reflect.Type) pretty.ValueWriter {
	var res pretty.ValueWriter
	if messageImplementsCache.ImplementedBy(typ) {
		res = vw
	}
	return res
}

// EnumValue represents an enum value.
type EnumValue struct {
	Number int32
	Name   string
}
