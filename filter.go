package pretty

import (
	"reflect"
)

// FilterWriter is a [ValueWriter] that calls the [ValueWriter] if the filter returns true.
//
// It should be created with [NewFilterWriter].
type FilterWriter[VW ValueWriter] struct {
	ValueWriter VW
	// Filter filters types.
	// The value is handled if it returns true or if it is nil.
	Filter func(typ reflect.Type) bool
}

// NewFilterWriter creates a new [FilterWriter].
func NewFilterWriter[VW ValueWriter](vw VW, f func(typ reflect.Type) bool) *FilterWriter[VW] {
	return &FilterWriter[VW]{
		ValueWriter: vw,
		Filter:      f,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FilterWriter[VW]) WriteValue(st *State, v reflect.Value) bool {
	return (vw.Filter == nil || vw.Filter(v.Type())) && vw.ValueWriter.WriteValue(st, v)
}

// Supports implements [SupportChecker].
func (vw *FilterWriter[VW]) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if vw.Filter == nil || vw.Filter(typ) {
		res = supportsValueWriter(typ, vw.ValueWriter)
	}
	return res
}

// FilterTypes returns a new [FilterWriter] filter function that returns true if the type is in the given list.
func FilterTypes(typs ...reflect.Type) func(typ reflect.Type) bool {
	set := make(map[reflect.Type]struct{}, len(typs))
	for _, typ := range typs {
		set[typ] = struct{}{}
	}
	return func(typ reflect.Type) bool {
		_, ok := set[typ]
		return ok
	}
}
