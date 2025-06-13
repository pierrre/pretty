package pretty

import (
	"reflect"
	"slices"

	"github.com/pierrre/pretty/internal/write"
)

// RecursionWriter is a [ValueWriter] that prevents recursion.
//
// It should be created with [NewRecursionWriter].
type RecursionWriter struct {
	ValueWriter
	// ShowInfos shows the infos (type and address).
	// Default: true.
	ShowInfos bool
}

// NewRecursionWriter creates a new [RecursionWriter].
func NewRecursionWriter(vw ValueWriter) *RecursionWriter {
	return &RecursionWriter{
		ValueWriter: vw,
		ShowInfos:   true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *RecursionWriter) WriteValue(st *State, v reflect.Value) bool {
	visitedAdded, recursionDetected := checkRecursion(st, v, vw.ShowInfos)
	if recursionDetected {
		return true
	}
	ok := vw.ValueWriter.WriteValue(st, v)
	if visitedAdded {
		postRecursion(st)
	}
	return ok
}

func checkRecursion(st *State, v reflect.Value, showInfos bool) (visitedAdded bool, recursionDetected bool) {
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
	write.MustString(st.Writer, "<recursion>")
	if showInfos {
		write.MustString(st.Writer, " ")
		writeType(st.Writer, e.Type)
		write.MustString(st.Writer, " ")
		writeUintptr(st.Writer, e.Addr)
	}
	return false, true
}

func postRecursion(st *State) {
	i := len(st.Visited) - 1
	st.Visited[i] = VisitedEntry{}
	st.Visited = st.Visited[:i]
}
