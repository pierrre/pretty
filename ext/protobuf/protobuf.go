// Package protobuf provides a [pretty.ValueWriter] for protobuf messages.
package protobuf

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty"
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
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	pm := v.Interface().(protoreflect.ProtoMessage) //nolint:forcetypeassert // Checked above.
	m := pm.ProtoReflect()
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
		must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(vw.getValueInterface(m.Get(fd), fd))))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	if hasFields {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

func (vw *ValueWriter) getValueInterface(v protoreflect.Value, fd protoreflect.FieldDescriptor) any {
	itf := v.Interface()
	switch itf := itf.(type) {
	case protoreflect.Message:
		return itf.Interface()
	case protoreflect.List:
		return vw.getValueInterfaceList(itf, fd)
	case protoreflect.Map:
		return vw.getValueInterfaceMap(itf, fd)
	case protoreflect.EnumNumber:
		return vw.getValueInterfaceEnum(itf, fd)
	}
	return itf
}

func (vw *ValueWriter) getValueInterfaceList(l protoreflect.List, fd protoreflect.FieldDescriptor) any {
	// TODO create typed slice
	res := make([]any, l.Len())
	for i := range l.Len() {
		res[i] = vw.getValueInterface(l.Get(i), fd)
	}
	return res
}

func (vw *ValueWriter) getValueInterfaceMap(m protoreflect.Map, fd protoreflect.FieldDescriptor) any {
	// TODO create typed map
	res := make(map[any]any, m.Len())
	for key, value := range m.Range {
		res[vw.getValueInterface(key.Value(), fd.MapKey())] = vw.getValueInterface(value, fd.MapValue())
	}
	return res
}

func (vw *ValueWriter) getValueInterfaceEnum(e protoreflect.EnumNumber, fd protoreflect.FieldDescriptor) EnumValue {
	res := EnumValue{
		Number: int32(e),
	}
	ed := fd.Enum().Values().ByNumber(e)
	if ed != nil {
		res.Name = string(ed.Name())
	}
	return res
}

// EnumValue represents an enum value.
type EnumValue struct {
	Number int32
	Name   string
}
