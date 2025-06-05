package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
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
func (vw *IntValueWriter) WriteValue(st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return false
	}
	write.Must(strconvio.WriteInt(st.Writer, v.Int(), vw.Base))
	return true
}

// Supports implements [SupportChecker].
func (vw *IntValueWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	switch typ.Kind() { //nolint:exhaustive // Only handles int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res = vw
	}
	return res
}
