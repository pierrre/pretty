package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// StructValueWriter is a [ValueWriter] that handles struct values.
//
// It should be created with [NewStructValueWriter].
type StructValueWriter struct {
	ValueWriter
	// FieldFilter filters the fields.
	// Default: nil.
	FieldFilter StructFieldFilter
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
	write.MustString(st.Writer, "{")
	fields := reflectutil.GetStructFields(v.Type())
	hasFields := false
	st.IndentLevel++
	for i, field := range fields.Range {
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
		st.KnownType = false // We want to show the types of fields and values.
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
func (vw *StructValueWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Struct {
		res = vw
	}
	return res
}

// StructFieldFilter is a function that filters struct fields.
//
// It's used by [StructValueWriter].
type StructFieldFilter func(v reflect.Value, field reflect.StructField) bool

// NewExportedStructFieldFilter creates a new [StructFieldFilter] that returns true for exported fields and false otherwise.
func NewExportedStructFieldFilter() StructFieldFilter {
	return exportedStructFieldFilter
}

func exportedStructFieldFilter(v reflect.Value, field reflect.StructField) bool { //nolint:gocritic // The StructField type is large, but we need to use it.
	return field.IsExported()
}
