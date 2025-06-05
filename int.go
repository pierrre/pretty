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
	if !vw.supportsType(v.Type()) {
		return false
	}
	return vw.writeValue(st, v)
}

func (vw *IntValueWriter) writeValue(st *State, v reflect.Value) bool {
	write.Must(strconvio.WriteInt(st.Writer, v.Int(), vw.Base))
	return true
}

// SupportsType implements [TypeSupportChecker].
func (vw *IntValueWriter) SupportsType(typ reflect.Type) ValueWriterFunc {
	if vw.supportsType(typ) {
		return vw.writeValue
	}
	return nil
}

func (vw *IntValueWriter) supportsType(typ reflect.Type) bool {
	switch typ.Kind() { //nolint:exhaustive // Only handles int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}
