package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal"
)

// StructValueWriter is a [ValueWriter] that handles struct values.
//
// It should be created with [NewStructValueWriter].
type StructValueWriter struct {
	ValueWriter
	// FieldFilter filters the fields.
	// Default: nil.
	FieldFilter func(v reflect.Value, field reflect.StructField) bool
}

// NewStructValueWriter creates a new [StructValueWriter] with default values.
func NewStructValueWriter(vw ValueWriter) *StructValueWriter {
	return &StructValueWriter{
		ValueWriter: vw,
		FieldFilter: nil,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StructValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct {
		return false
	}
	defer st.SetRestoreKnownType(false)() // We want to show the types of fields and values.
	internal.MustWriteString(st.Writer, "{")
	fields := getStructFields(v.Type())
	hasFields := false
	st.IndentLevel++
	for i, field := range fields {
		if vw.FieldFilter != nil && !vw.FieldFilter(v, field) {
			continue
		}
		if !hasFields {
			internal.MustWriteString(st.Writer, "\n")
			hasFields = true
		}
		st.writeIndent()
		internal.MustWriteString(st.Writer, field.Name)
		internal.MustWriteString(st.Writer, ": ")
		internal.MustHandle(vw.ValueWriter(st, v.Field(i)))
		internal.MustWriteString(st.Writer, ",\n")
	}
	st.IndentLevel--
	if hasFields {
		st.writeIndent()
	}
	internal.MustWriteString(st.Writer, "}")
	return true
}

var structFieldsCache syncutil.Map[reflect.Type, []reflect.StructField]

func getStructFields(typ reflect.Type) []reflect.StructField {
	fields, ok := structFieldsCache.Load(typ)
	if ok {
		return fields
	}
	for i := range typ.NumField() {
		field := typ.Field(i)
		fields = append(fields, field)
	}
	structFieldsCache.Store(typ, fields)
	return fields
}

// UnsafePointerValueWriter is a struct field filter that returns true for exported fields.
func ExportedStructFieldFilter(v reflect.Value, field reflect.StructField) bool { //nolint:gocritic // The StructField type is large, but we need to use it.
	return field.IsExported()
}
