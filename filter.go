package pretty

import (
	"reflect"
)

// FilterWriter is a [ValueWriter] that calls the [ValueWriter] if the filter returns true.
//
// It should be created with [NewFilterWriter].
type FilterWriter[VW ValueWriter] struct {
	ValueWriter VW
	// Filter filters values.
	// The value is handled if it returns true or if it is nil.
	Filter func(v reflect.Value) bool
}

// NewFilterWriter creates a new [FilterWriter].
func NewFilterWriter[VW ValueWriter](vw VW, f func(v reflect.Value) bool) *FilterWriter[VW] {
	return &FilterWriter[VW]{
		ValueWriter: vw,
		Filter:      f,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FilterWriter[VW]) WriteValue(st *State, v reflect.Value) bool {
	return (vw.Filter == nil || vw.Filter(v)) && vw.ValueWriter.WriteValue(st, v)
}
