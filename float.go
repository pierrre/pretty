package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// FloatWriter is a [ValueWriter] that handles float values.
//
// It should be created with [NewFloatWriter].
type FloatWriter struct {
	// Format is the format used to format the float.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the float.
	// Default: -1.
	Precision int
}

// NewFloatWriter creates a new [FloatWriter] with default values.
func NewFloatWriter() *FloatWriter {
	return &FloatWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FloatWriter) WriteValue(st *State, v reflect.Value) bool {
	var bitSize int
	switch v.Kind() { //nolint:exhaustive // Only handles float.
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	default:
		return false
	}
	write.Must(strconvio.WriteFloat(st.Writer, v.Float(), vw.Format, vw.Precision, bitSize))
	return true
}

// Supports implements [SupportChecker].
func (vw *FloatWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	switch typ.Kind() { //nolint:exhaustive // Only handles float.
	case reflect.Float32, reflect.Float64:
		res = vw
	}
	return res
}
