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

// ConfigureDefault configures [pretty.DefaultWriter] with [ConfigureCommonWriterDefault].
func ConfigureDefault() {
	ConfigureCommonWriterDefault(pretty.DefaultWriter)
}

// ConfigureCommonWriterDefault configures a [pretty.CommonWriter] with a default [MessageWriter].
func ConfigureCommonWriterDefault(vw *pretty.CommonWriter) {
	ConfigureCommonWriter(vw, NewMessageWriter(vw))
}

// ConfigureCommonWriter configures a [pretty.CommonWriter] with a [MessageWriter].
func ConfigureCommonWriter(vw *pretty.CommonWriter, mw *MessageWriter) {
	vw.ValueWriters = append(vw.ValueWriters, mw)
}

// MessageWriter is a [pretty.MessageWriter] that handles protobuf messages.
//
// It should be created with [NewMessageWriter].
type MessageWriter struct {
	pretty.ValueWriter
	// ShowFieldsType shows the type of the fields.
	// Default: true.
	ShowFieldsType bool
}

// NewMessageWriter creates a new [MessageWriter].
func NewMessageWriter(vw pretty.ValueWriter) *MessageWriter {
	return &MessageWriter{
		ValueWriter:    vw,
		ShowFieldsType: true,
	}
}

// WriteValue implements [pretty.ValueWriter].
func (vw *MessageWriter) WriteValue(st *pretty.State, v reflect.Value) bool {
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

func (vw *MessageWriter) writeMessage(st *pretty.State, m protoreflect.Message) {
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
		st.KnownType = !vw.ShowFieldsType
		must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(vw.getInterface(m.Get(fd), fd))))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	if hasFields {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
}

func (vw *MessageWriter) getInterface(v protoreflect.Value, fd protoreflect.FieldDescriptor) any {
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

func (vw *MessageWriter) getList(l protoreflect.List, fd protoreflect.FieldDescriptor) any {
	// TODO create typed slice
	res := make([]any, l.Len())
	for i := range l.Len() {
		res[i] = vw.getInterface(l.Get(i), fd)
	}
	return res
}

func (vw *MessageWriter) getMap(m protoreflect.Map, fd protoreflect.FieldDescriptor) any {
	// TODO create typed map
	res := make(map[any]any, m.Len())
	for key, value := range m.Range {
		res[vw.getInterface(key.Value(), fd.MapKey())] = vw.getInterface(value, fd.MapValue())
	}
	return res
}

func (vw *MessageWriter) getEnum(e protoreflect.EnumNumber, fd protoreflect.FieldDescriptor) EnumValue {
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
func (vw *MessageWriter) Supports(typ reflect.Type) pretty.ValueWriter {
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
