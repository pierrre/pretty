package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// BoolValueWriter is a [ValueWriter] that handles bool values.
//
// It should be created with [NewBoolValueWriter].
type BoolValueWriter struct{}

// NewBoolValueWriter creates a new [BoolValueWriter].
func NewBoolValueWriter() *BoolValueWriter {
	return &BoolValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *BoolValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Bool {
		return false
	}
	write.Must(strconvio.WriteBool(st.Writer, v.Bool()))
	return true
}
