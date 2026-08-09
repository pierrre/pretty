package pretty

import (
	"reflect"
	"testing"
)

func TestStateReleaseLargeVisited(t *testing.T) {
	st := newState("")
	vw := NewRecursionWriter(ValueWriterFunc(func(*State, reflect.Value) bool { return true }))
	ptrs := make([]*int, stateVisitedKeepMax+4)
	for i := range ptrs {
		ptrs[i] = new(int)
		vw.checkRecursion(st, reflect.ValueOf(ptrs[i]))
	}
	st.release()
	if st.Visited != nil {
		t.Fatalf("state Visited map not released for large graph: %v", st.Visited)
	}
}

func TestStateReleaseSmallVisited(t *testing.T) {
	st := newState("")
	vw := NewRecursionWriter(ValueWriterFunc(func(*State, reflect.Value) bool { return true }))
	ptrs := make([]*int, stateVisitedKeepMax)
	for i := range ptrs {
		ptrs[i] = new(int)
		vw.checkRecursion(st, reflect.ValueOf(ptrs[i]))
	}
	st.release()
	if st.Visited == nil {
		t.Fatal("state Visited map released for small graph")
	}
}
