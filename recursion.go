package pretty

import (
	"reflect"
	"slices"

	"github.com/pierrre/go-libs/reflectutil"
)

// RecursionWriter is a [ValueWriter] that prevents recursion.
//
// It should be created with [NewRecursionWriter].
type RecursionWriter struct {
	ValueWriter
	// ShowAddr shows the address (and type).
	// Default: true.
	ShowAddr bool
}

// NewRecursionWriter creates a new [RecursionWriter].
func NewRecursionWriter(vw ValueWriter) *RecursionWriter {
	return &RecursionWriter{
		ValueWriter: vw,
		ShowAddr:    true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *RecursionWriter) WriteValue(st *State, v reflect.Value) bool {
	visitedAdded, recursionDetected := vw.checkRecursion(st, v)
	if recursionDetected {
		return true
	}
	ok := vw.ValueWriter.WriteValue(st, v)
	if visitedAdded {
		vw.postRecursion(st)
	}
	return ok
}

func (vw *RecursionWriter) checkRecursion(st *State, v reflect.Value) (visitedAdded bool, recursionDetected bool) {
	switch v.Kind() { //nolint:exhaustive // Only handles pointer kinds.
	case reflect.Pointer, reflect.Map, reflect.Slice:
	default:
		return false, false
	}
	e := VisitedEntry{
		Type: v.Type(),
		Addr: uintptr(v.UnsafePointer()),
	}
	if !slices.Contains(st.Visited, e) {
		st.Visited = append(st.Visited, e)
		return true, false
	}
	st.Writer.AppendString("<recursion>")
	if vw.ShowAddr {
		st.Writer.AppendByte(' ')
		st.Writer.AppendString(reflectutil.TypeFullName(e.Type))
		st.Writer.AppendByte(' ')
		writeUintptr(st, e.Addr)
	}
	return false, true
}

func (vw *RecursionWriter) postRecursion(st *State) {
	i := len(st.Visited) - 1
	st.Visited[i] = VisitedEntry{}
	st.Visited = st.Visited[:i]
}
