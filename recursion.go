package pretty

import (
	"reflect"
	"slices"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/write"
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
	switch v.Kind() {
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
	write.MustString(st.Writer, "<recursion>")
	if vw.ShowAddr {
		write.MustString(st.Writer, " ")
		write.MustString(st.Writer, reflectutil.TypeFullName(e.Type))
		write.MustString(st.Writer, " ")
		writeUintptr(st.Writer, e.Addr)
	}
	return false, true
}

func (vw *RecursionWriter) postRecursion(st *State) {
	i := len(st.Visited) - 1
	st.Visited[i] = VisitedEntry{}
	st.Visited = st.Visited[:i]
}
