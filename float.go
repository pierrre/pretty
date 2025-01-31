package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// FloatValueWriter is a [ValueWriter] that handles float values.
//
// It should be created with [NewFloatValueWriter].
type FloatValueWriter struct {
	// Format is the format used to format the float.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the float.
	// Default: -1.
	Precision int
}

// NewFloatValueWriter creates a new [FloatValueWriter] with default values.
func NewFloatValueWriter() *FloatValueWriter {
	return &FloatValueWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FloatValueWriter) WriteValue(st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles float.
	case reflect.Float32, reflect.Float64:
		write.Must(strconvio.WriteFloat(st.Writer, v.Float(), vw.Format, vw.Precision, v.Type().Bits()))
		return true
	}
	return false
}
