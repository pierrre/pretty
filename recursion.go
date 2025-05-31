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
	visitedAdded, recursionDetected := checkRecursion(st, v)
	if recursionDetected {
		return true
	}
	ok := vw.ValueWriter.WriteValue(st, v)
	if visitedAdded {
		postRecursion(st)
	}
	return ok
}

func checkRecursion(st *State, v reflect.Value) (visitedAdded bool, recursionDetected bool) {
	switch v.Kind() { //nolint:exhaustive // Only handles pointer kinds.
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
	default:
		return false, false
	}
	vp := v.Pointer()
	if slices.Contains(st.Visited, vp) {
		write.MustString(st.Writer, "<recursion>")
		return false, true
	}
	st.Visited = append(st.Visited, vp)
	return true, false
}

func postRecursion(st *State) {
	st.Visited = st.Visited[:len(st.Visited)-1]
}
