package pretty

import (
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
func (vw *FilterValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !vw.Filter(v) {
		return false
	}
	return vw.ValueWriter.WriteValue(st, v)
}
