package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// StructWriter is a [ValueWriter] that handles struct values.
//
// It should be created with [NewStructWriter].
type StructWriter struct {
	ValueWriter
	// FieldFilter filters the fields.
	// Default: nil.
	FieldFilter StructFieldFilter
	// ShowFieldsType shows the type of the fields.
	// Default: true.
	ShowFieldsType bool
}

// NewStructWriter creates a new [StructWriter] with default values.
func NewStructWriter(vw ValueWriter) *StructWriter {
	return &StructWriter{
		ValueWriter:    vw,
		FieldFilter:    nil,
		ShowFieldsType: true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StructWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct {
		return false
	}
	write.MustString(st.Writer, "{")
	fields := reflectutil.GetStructFields(v.Type())
	hasFields := false
	st.IndentLevel++
	for i := range fields.Len() {
		field := fields.Get(i)
		if vw.FieldFilter != nil && !vw.FieldFilter(v, field) {
			continue
		}
		if !hasFields {
			write.MustString(st.Writer, "\n")
			hasFields = true
		}
		st.WriteIndent()
		write.MustString(st.Writer, field.Name)
		write.MustString(st.Writer, ": ")
		st.KnownType = !vw.ShowFieldsType
		must.Handle(vw.ValueWriter.WriteValue(st, v.Field(i)))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	if hasFields {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

// Supports implements [SupportChecker].
func (vw *StructWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Struct {
		res = vw
	}
	return res
}

// StructFieldFilter is a function that filters struct fields.
//
// It's used by [StructWriter].
type StructFieldFilter func(v reflect.Value, field reflect.StructField) bool

// NewExportedStructFieldFilter creates a new [StructFieldFilter] that returns true for exported fields and false otherwise.
func NewExportedStructFieldFilter() StructFieldFilter {
	return exportedStructFieldFilter
}

func exportedStructFieldFilter(v reflect.Value, field reflect.StructField) bool { //nolint:gocritic // The StructField type is large, but we need to use it.
	return field.IsExported()
}
