package pretty

import (
	"io"
	"reflect"
	"strconv"
)

// ComplexValueWriter is a [ValueWriter] that handles complex values.
//
// It should be created with [NewComplexValueWriter].
type ComplexValueWriter struct {
	// Format is the format used to format the complex.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the complex.
	// Default: -1.
	Precision int
}

// NewComplexValueWriter creates a new [ComplexValueWriter] with default values.
func NewComplexValueWriter() *ComplexValueWriter {
	return &ComplexValueWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ComplexValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles complex.
	case reflect.Complex64, reflect.Complex128:
		writeString(w, strconv.FormatComplex(v.Complex(), vw.Format, vw.Precision, v.Type().Bits()))
		return true
	}
	return false
}