package pretty

import (
	"io"

	"github.com/pierrre/go-libs/syncutil"
)

// State represents the state of the [Printer].
//
// Functions must restore the original state when they return.
type State struct {
	Depth        int
	IndentString string
	IndentLevel  int
	Visited      *[]uintptr
	KnownType    bool
}

func (st State) writeIndent(w io.Writer) {
	writeIndent(w, st.IndentString, st.IndentLevel)
}

func getState() State {
	vs := stateVisitedPool.Get()
	*vs = (*vs)[:0]
	return State{
		Visited: vs,
	}
}

func (st State) pushPopVisited(p uintptr) func() {
	st.pushVisited(p)
	return st.popVisited
}

func (st State) pushVisited(p uintptr) {
	*st.Visited = append(*st.Visited, p)
}

func (st State) popVisited() {
	s := *st.Visited
	*st.Visited = s[:len(s)-1]
}

func (st State) release() {
	stateVisitedPool.Put(st.Visited)
}

var stateVisitedPool = syncutil.PoolFor[*[]uintptr]{
	New: func() *[]uintptr {
		return new([]uintptr)
	},
}
