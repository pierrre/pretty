package pretty

import (
	"io"
	"reflect"
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
func (vw *MaxDepthValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if vw.Max <= 0 {
		return vw.ValueWriter(w, st, v)
	}
	if st.Depth >= vw.Max {
		writeString(w, "<max depth>")
		return true
	}
	st.Depth++
	return vw.ValueWriter(w, st, v)
}
