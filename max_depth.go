package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/write"
)

// MaxDepthWriter is a [ValueWriter] that limits the depth.
//
// It should be created with [NewMaxDepthWriter].
type MaxDepthWriter struct {
	ValueWriter
	// Max is the maximum depth.
	// Default: 0 (no limit).
	Max int
}

// NewMaxDepthWriter creates a new [MaxDepthWriter].
func NewMaxDepthWriter(vw ValueWriter) *MaxDepthWriter {
	return &MaxDepthWriter{
		ValueWriter: vw,
		Max:         0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MaxDepthWriter) WriteValue(st *State, v reflect.Value) bool {
	maxReached := vw.checkMaxDepth(st)
	if maxReached {
		return true
	}
	ok := vw.ValueWriter.WriteValue(st, v)
	vw.postMaxDepth(st)
	return ok
}

func (vw *MaxDepthWriter) checkMaxDepth(st *State) (maxReached bool) {
	if vw.Max > 0 && st.Depth >= vw.Max {
		write.MustString(st.Writer, "<max depth>")
		maxReached = true
	}
	st.Depth++
	return maxReached
}

func (vw *MaxDepthWriter) postMaxDepth(st *State) {
	st.Depth--
}
