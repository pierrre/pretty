package pretty

import (
	"reflect"
	"slices"

	"github.com/pierrre/pretty/internal/write"
)

// RecursionValueWriter is a [ValueWriter] that prevents recursion.
//
// It should be created with [NewRecursionValueWriter].
type RecursionValueWriter struct {
	ValueWriter
}

// NewRecursionValueWriter creates a new [RecursionValueWriter].
func NewRecursionValueWriter(vw ValueWriter) *RecursionValueWriter {
	return &RecursionValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *RecursionValueWriter) WriteValue(st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles pointer kinds.
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
	default:
		return vw.ValueWriter(st, v)
	}
	vp := v.Pointer()
	if slices.Contains(st.Visited, vp) {
		write.MustString(st.Writer, "<recursion>")
		return true
	}
	st.Visited = append(st.Visited, vp)
	defer func() {
		st.Visited = st.Visited[:len(st.Visited)-1]
	}()
	return vw.ValueWriter(st, v)
}
