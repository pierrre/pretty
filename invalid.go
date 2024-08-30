package pretty

import (
	"io"
	"reflect"
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
	return checkInvalid(st.Writer, v)
}

func checkInvalid(w io.Writer, v reflect.Value) bool {
	if v.IsValid() {
		return false
	}
	writeString(w, "<invalid>")
	return true
}
