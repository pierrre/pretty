package pretty

import (
	"io"
	"reflect"
	"slices"
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
func (vw *RecursionValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles pointer kinds.
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
	default:
		return vw.ValueWriter(w, st, v)
	}
	vp := v.Pointer()
	if slices.Contains(*st.Visited, vp) {
		writeString(w, "<recursion>")
		return true
	}
	defer st.pushPopVisited(vp)()
	return vw.ValueWriter(w, st, v)
}
