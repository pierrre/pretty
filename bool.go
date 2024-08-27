package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
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
func (vw *BoolValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() != reflect.Bool {
		return false
	}
	mustWrite(strconvio.WriteBool(w, v.Bool()))
	return true
}
