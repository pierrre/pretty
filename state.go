package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/bytesutil"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/indent"
)

// State represents the state of the [Printer].
//
// Functions must restore the original state when they return.
type State struct {
	Writer       bytesutil.Writer
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

func newState(indentString string) *State {
	st := statePool.Get()
	st.Depth = 0
	st.IndentString = indentString
	st.IndentLevel = 0
	st.Visited = st.Visited[:0]
	st.KnownType = false
	st.ShowInfos = true
	return st
}

// WriteIndent writes the current indentation to the writer.
func (st *State) WriteIndent() {
	st.Writer = indent.Append(st.Writer, st.IndentString, st.IndentLevel)
}

func (st *State) release() {
	st.Writer.Reset()
	statePool.Put(st)
}

// VisitedEntry represents a visited value.
type VisitedEntry struct {
	Type reflect.Type
	Addr uintptr
}
