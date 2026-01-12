package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// IntWriter is a [ValueWriter] that handles int values.
//
// It should be created with [NewIntWriter].
type IntWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewIntWriter creates a new [IntWriter] with default values.
func NewIntWriter() *IntWriter {
	return &IntWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IntWriter) WriteValue(st *State, v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return false
	}
	write.Must(strconvio.WriteInt(st.Writer, v.Int(), vw.Base))
	return true
}

// Supports implements [SupportChecker].
func (vw *IntWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res = vw
	}
	return res
}
