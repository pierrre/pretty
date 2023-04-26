// Package pierrreerrors provides an integration with github.com/pierrre/errors.
package pierrreerrors

import (
	"io"
	"reflect"

	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/pretty"
)

// ConfigureDefault calls Configure() with pretty.DefaultConfig.
func ConfigureDefault() {
	Configure(pretty.DefaultConfig)
}

// Configure configures a pretty.Config with NewValueWriter().
// It prepends the ValueWriter at the beginning of the ValueWriters slice.
func Configure(c *pretty.Config) {
	c.ValueWriters = append([]pretty.ValueWriter{NewValueWriter()}, c.ValueWriters...)
}

// NewValueWriter creates a new pretty.ValueWriter for github.com/pierrre/errors.
//
// It prints the verbose representation of the error.
func NewValueWriter() pretty.ValueWriter {
	return write
}

var typeError = reflect.TypeOf((*error)(nil)).Elem()

func write(c *pretty.Config, w io.Writer, st *pretty.State, v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	if !v.Type().Implements(typeError) {
		return false
	}
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	st.Indent++
	iw := pretty.GetIndentWriter(w, c, st, true)
	errverbose.Write(iw, err)
	iw.Release()
	st.Indent--
	return true
}
