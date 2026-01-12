package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// ComplexWriter is a [ValueWriter] that handles complex values.
//
// It should be created with [NewComplexWriter].
type ComplexWriter struct {
	// Format is the format used to format the complex.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the complex.
	// Default: -1.
	Precision int
}

// NewComplexWriter creates a new [ComplexWriter] with default values.
func NewComplexWriter() *ComplexWriter {
	return &ComplexWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ComplexWriter) WriteValue(st *State, v reflect.Value) bool {
	var bitSize int
	switch v.Kind() {
	case reflect.Complex64:
		bitSize = 64
	case reflect.Complex128:
		bitSize = 128
	default:
		return false
	}
	write.Must(strconvio.WriteComplex(st.Writer, v.Complex(), vw.Format, vw.Precision, bitSize))
	return true
}

// Supports implements [SupportChecker].
func (vw *ComplexWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	switch typ.Kind() {
	case reflect.Complex64, reflect.Complex128:
		res = vw
	}
	return res
}
