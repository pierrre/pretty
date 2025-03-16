package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/pretty/internal/write"
)

// InvalidValueWriter is a [ValueWriter] that handles invalid values.
//
// It should be created with [NewInvalidValueWriter].
type InvalidValueWriter struct{}

// NewInvalidValueWriter creates a new [InvalidValueWriter].
func NewInvalidValueWriter() *InvalidValueWriter {
	return &InvalidValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *InvalidValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.IsValid() {
		return false
	}
	write.MustString(st.Writer, "<invalid>")
	return true
}

func checkInvalidNil(w io.Writer, v reflect.Value) bool {
	if v.IsValid() {
		return false
	}
	writeNil(w)
	return true
}
