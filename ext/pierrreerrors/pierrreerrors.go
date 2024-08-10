// Package pierrreerrors provides an integration with github.com/pierrre/errors.
package pierrreerrors

import (
	"io"
	"reflect"
	"slices"

	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/pretty"
)

// ConfigureDefault configures [pretty.DefaultCommonValueWriter] with [DefaultValueWriter].
//
// It calls [ConfigureCommonValueWriter] with [pretty.DefaultCommonValueWriter].
func ConfigureDefault() {
	ConfigureCommonValueWriter(pretty.DefaultCommonValueWriter)
}

// Configure configures a [pretty.CommonValueWriter] with [DefaultValueWriter].
//
// It calls [ConfigureValueWriters] with [pretty.CommonValueWriter.ValueWriters].
func ConfigureCommonValueWriter(vw *pretty.CommonValueWriter) {
	vw.ValueWriters = ConfigureValueWriters(vw.ValueWriters)
}

// ConfigureValueWriters configures a [pretty.ValueWriters] with [DefaultValueWriter].
//
// It prepends the [ValueWriter.WriteValue] at the beginning of the slice.
func ConfigureValueWriters(vws pretty.ValueWriters) pretty.ValueWriters {
	return slices.Insert(vws, 0, NewValueWriter().WriteValue)
}

// DefaultValueWriter is the default [ValueWriter].
var DefaultValueWriter = NewValueWriter()

// ValueWriter is a [pretty.ValueWriter] that handles errors and write them with [errverbose.Write].
//
// It writes the verbose representation of the error.
//
// It should be created with [NewValueWriter].
type ValueWriter struct{}

// NewValueWriter creates a new [ValueWriter].
func NewValueWriter() *ValueWriter {
	return &ValueWriter{}
}

var typeError = reflect.TypeFor[error]()

// WriteValue implements [pretty.ValueWriter].
func (vw *ValueWriter) WriteValue(c *pretty.Config, w io.Writer, st *pretty.State, v reflect.Value) bool {
	if !v.Type().Implements(typeError) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	st.RunIndent(func(st *pretty.State) {
		iw := pretty.GetIndentWriter(c, w, st, true)
		defer iw.Release()
		errverbose.Write(iw, err)
	})
	return true
}
