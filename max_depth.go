package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/write"
)

// MaxDepthValueWriter is a [ValueWriter] that limits the depth.
//
// It should be created with [NewMaxDepthValueWriter].
type MaxDepthValueWriter struct {
	ValueWriter
	// Max is the maximum depth.
	// Default: 0 (no limit).
	Max int
}

// NewMaxDepthValueWriter creates a new [MaxDepthValueWriter].
func NewMaxDepthValueWriter(vw ValueWriter) *MaxDepthValueWriter {
	return &MaxDepthValueWriter{
		ValueWriter: vw,
		Max:         0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MaxDepthValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.Max <= 0 {
		return vw.ValueWriter.WriteValue(st, v)
	}
	if checkMaxDepth(st, vw.Max) {
		return true
	}
	ok := vw.ValueWriter.WriteValue(st, v)
	postMaxDepth(st)
	return ok
}

func checkMaxDepth(st *State, maxDepth int) bool {
	if st.Depth >= maxDepth {
		write.MustString(st.Writer, "<max depth>")
		return true
	}
	st.Depth++
	return false
}

func postMaxDepth(st *State) {
	st.Depth--
}
