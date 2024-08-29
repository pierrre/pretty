package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal"
)

// IntValueWriter is a [ValueWriter] that handles int values.
//
// It should be created with [NewIntValueWriter].
type IntValueWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewIntValueWriter creates a new [IntValueWriter] with default values.
func NewIntValueWriter() *IntValueWriter {
	return &IntValueWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IntValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		internal.MustWrite(strconvio.WriteInt(w, v.Int(), vw.Base))
		return true
	}
	return false
}
