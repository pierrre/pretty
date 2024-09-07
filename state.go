package pretty

import (
	"io"

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
	Visited      []uintptr
	KnownType    bool
	ShowInfos    bool
}

var statePool = syncutil.PoolFor[*State]{
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

func (st *State) writeIndent() {
	indent.MustWrite(st.Writer, st.IndentString, st.IndentLevel)
}

func (st *State) setRestoreKnownType(knownType bool) func() {
	st.KnownType, knownType = knownType, st.KnownType
	return func() {
		st.KnownType = knownType
	}
}

func (st *State) release() {
	st.Writer = nil
	statePool.Put(st)
}
