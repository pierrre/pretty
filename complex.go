package pretty

import (
	"reflect"
	"strconv"
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
	switch v.Kind() { //nolint:exhaustive // Only handles complex.
	case reflect.Complex64:
		bitSize = 64
	case reflect.Complex128:
		bitSize = 128
	default:
		return false
	}
	st.Writer = appendComplex(st.Writer, v.Complex(), vw.Format, vw.Precision, bitSize)
	return true
}

// Supports implements [SupportChecker].
func (vw *ComplexWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	switch typ.Kind() { //nolint:exhaustive // Only handles complex.
	case reflect.Complex64, reflect.Complex128:
		res = vw
	}
	return res
}

func appendComplex(dst []byte, c complex128, fmt byte, prec, bitSize int) []byte {
	bitSize >>= 1 // complex64 uses float32 internally
	dst = append(dst, '(')
	dst = strconv.AppendFloat(dst, real(c), fmt, prec, bitSize)
	i := len(dst)
	dst = strconv.AppendFloat(dst, imag(c), fmt, prec, bitSize)
	// Check if imaginary part has a sign. If not, add one.
	if dst[i] != '+' && dst[i] != '-' {
		dst = append(dst, 0)
		copy(dst[i+1:], dst[i:])
		dst[i] = '+'
	}
	dst = append(dst, "i)"...)
	return dst
}
