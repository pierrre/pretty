package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
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

// FilterTypes returns a new [FilterWriter] filter function that returns true if the type is in the given list or if it implements any of the given interface types.
func FilterTypes(typs ...reflect.Type) func(typ reflect.Type) bool {
	set := make(map[reflect.Type]struct{}, len(typs))
	var ics []*reflectutil.ImplementsCache
	for _, typ := range typs {
		if _, ok := set[typ]; !ok {
			set[typ] = struct{}{}
			if typ.Kind() == reflect.Interface {
				ics = append(ics, reflectutil.NewImplementsCache(typ))
			}
		}
	}
	return func(typ reflect.Type) bool {
		_, ok := set[typ]
		if ok {
			return true
		}
		for _, ic := range ics {
			if ic.ImplementedBy(typ) {
				return true
			}
		}
		return false
	}
}
