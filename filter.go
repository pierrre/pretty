package pretty

import (
	"io"
	"reflect"
)

// FilterValueWriter is a [ValueWriter] that calls the [ValueWriter] if the filter returns true.
//
// It should be created with [NewFilterValueWriter].
type FilterValueWriter struct {
	ValueWriter
	Filter func(v reflect.Value) bool
}

// NewFilterValueWriter creates a new [FilterValueWriter].
func NewFilterValueWriter(vw ValueWriter, f func(v reflect.Value) bool) *FilterValueWriter {
	return &FilterValueWriter{
		ValueWriter: vw,
		Filter:      f,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FilterValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if !vw.Filter(v) {
		return false
	}
	return vw.ValueWriter(w, st, v)
}
