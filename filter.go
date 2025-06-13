package pretty

import (
	"reflect"
)

// FilterWriter is a [ValueWriter] that calls the [ValueWriter] if the filter returns true.
//
// It should be created with [NewFilterWriter].
type FilterWriter struct {
	ValueWriter
	Filter func(v reflect.Value) bool
}

// NewFilterWriter creates a new [FilterWriter].
func NewFilterWriter(vw ValueWriter, f func(v reflect.Value) bool) *FilterWriter {
	return &FilterWriter{
		ValueWriter: vw,
		Filter:      f,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FilterWriter) WriteValue(st *State, v reflect.Value) bool {
	if !vw.Filter(v) {
		return false
	}
	return vw.ValueWriter.WriteValue(st, v)
}
