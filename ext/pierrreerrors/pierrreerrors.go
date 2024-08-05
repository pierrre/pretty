// Package pierrreerrors provides an integration with github.com/pierrre/errors.
package pierrreerrors

import (
	"io"
	"reflect"
	"slices"

	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/pretty"
)

// ConfigureDefault calls [Configure] with [pretty.DefaultConfig].
func ConfigureDefault() {
	Configure(pretty.DefaultConfig)
}

// Configure configures a [pretty.Config] with [NewValueWriter].
// It prepends the [pretty.ValueWriter] at the beginning of ValueWriter slice.
func Configure(c *pretty.Config) {
	c.ValueWriters = slices.Insert(c.ValueWriters, 0, NewValueWriter())
}

// NewValueWriter creates a new [pretty.ValueWriter] for github.com/pierrre/errors.
//
// It prints the verbose representation of the error.
func NewValueWriter() pretty.ValueWriter {
	return write
}

var typeError = reflect.TypeFor[error]()

func write(c *pretty.Config, w io.Writer, st *pretty.State, v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	if !v.Type().Implements(typeError) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	st.RunIndent(func(st *pretty.State) {
		iw := pretty.GetIndentWriter(w, c, st, true)
		defer iw.Release()
		errverbose.Write(iw, err)
	})
	return true
}
