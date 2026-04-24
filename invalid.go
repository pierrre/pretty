package pretty

import (
	"reflect"
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
	st.Writer.AppendString("<invalid>")
	return true
}

func checkInvalidNil(st *State, v reflect.Value) bool {
	if v.IsValid() {
		return false
	}
	writeNil(st)
	return true
}
