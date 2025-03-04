// Package protobuf provides a [pretty.ValueWriter] for protobuf messages.
package protobuf

import (
	"reflect"

	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var protoMessageType = reflect.TypeFor[protoreflect.ProtoMessage]()

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
	if !v.Type().Implements(protoMessageType) {
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
	defer st.SetRestoreKnownType(false)() // We want to show the types of fields and values.
	write.MustString(st.Writer, "{")
	fs := m.Descriptor().Fields()
	l := fs.Len()
	hasFields := false
	st.IndentLevel++
	for i := range l {
		fd := fs.Get(i)
		if !hasFields {
			write.MustString(st.Writer, "\n")
			hasFields = true
		}
		st.WriteIndent()
		write.MustString(st.Writer, string(fd.Name()))
		write.MustString(st.Writer, ": ")
		must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(vw.getValueInterface(m.Get(fd)))))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	if hasFields {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

func (vw *ValueWriter) getValueInterface(v protoreflect.Value) any {
	itf := v.Interface()
	switch itf := itf.(type) {
	case protoreflect.Message:
		return itf.Interface()
	case protoreflect.List:
		return vw.getValueInterfaceList(itf)
	case protoreflect.Map:
		return vw.getValueInterfaceMap(itf)
	}
	return itf
}

func (vw *ValueWriter) getValueInterfaceList(l protoreflect.List) any {
	res := make([]any, l.Len())
	for i := range l.Len() {
		res[i] = vw.getValueInterface(l.Get(i))
	}
	return res
}

func (vw *ValueWriter) getValueInterfaceMap(m protoreflect.Map) any {
	res := make(map[any]any, m.Len())
	m.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
		res[key.Interface()] = vw.getValueInterface(value)
		return true
	})
	return res
}
