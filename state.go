package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/indent"
)

// State represents the state of the [Printer].
//
// Functions must restore the original state when they return.
type State struct {
	Writer       io.Writer
	Depth        int
	IndentString string
	IndentLevel  int
	Visited      []VisitedEntry
	KnownType    bool
	ShowInfos    bool
}

var statePool = syncutil.Pool[*State]{
	New: func() *State {
		return new(State)
	},
}

func newState(w io.Writer, indentString string) *State {
	st := statePool.Get()
	st.Writer = w
	st.Depth = 0
	st.IndentString = indentString
	st.IndentLevel = 0
	st.Visited = st.Visited[:0]
	st.KnownType = false
	st.ShowInfos = true
	return st
}

// WriteIndent writes the current indentation to the [io.Writer].
func (st *State) WriteIndent() {
	indent.MustWrite(st.Writer, st.IndentString, st.IndentLevel)
}

func (st *State) release() {
	st.Writer = nil
	statePool.Put(st)
}

// VisitedEntry represents a visited value.
type VisitedEntry struct {
	Type reflect.Type
	Addr uintptr
}
