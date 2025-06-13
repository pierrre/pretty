package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/pretty/internal/write"
)

// InvalidWriter is a [ValueWriter] that handles invalid values.
//
// It should be created with [NewInvalidWriter].
type InvalidWriter struct{}

// NewInvalidWriter creates a new [InvalidWriter].
func NewInvalidWriter() *InvalidWriter {
	return &InvalidWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *InvalidWriter) WriteValue(st *State, v reflect.Value) bool {
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
