package pretty

import (
	"reflect"
	"unicode/utf8"

	"github.com/pierrre/go-libs/bytesutil"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/indent"
)

// State represents the state of the [Printer].
//
// Functions must restore the original state when they return.
//
// KnownType is an exception: [TypeWriter] saves and restores it, but other writers (e.g. [StructWriter], [InterfaceWriter], [ErrorWriter]) set it without restoring, as a hint to the next [TypeWriter] in the chain.
// This is intentional and relies on TypeWriter wrapping the inner writer to re-establish the correct value.
type State struct {
	Writer       bytesutil.Writer
	Depth        int
	IndentString string
	IndentLevel  int
	Visited      map[VisitedEntry]struct{}
	KnownType    bool
	ShowInfos    bool
	// MaxRunes is the maximum number of runes of the output.
	// If the output exceeds this limit, it is truncated.
	// Default: 0 (no limit).
	MaxRunes int
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
	clear(st.Visited)
	st.KnownType = false
	st.ShowInfos = true
	st.MaxRunes = 0
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

func (st *State) truncateMaxRunes() {
	if st.MaxRunes <= 0 {
		return
	}
	w, ok := truncateRunes(st.Writer, st.MaxRunes)
	if ok {
		st.Writer = w
		writeTruncated(st)
	}
}

func truncateRunes(b []byte, limit int) ([]byte, bool) {
	n := 0
	for i := 0; i < len(b); {
		_, size := utf8.DecodeRune(b[i:])
		if n == limit {
			return b[:i], true
		}
		i += size
		n++
	}
	return b, false
}
