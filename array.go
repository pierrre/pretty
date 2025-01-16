package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal"
	"github.com/pierrre/pretty/internal/must"
)

// ArrayValueWriter is a [ValueWriter] that handles array values.
//
// It should be created with [NewArrayValueWriter].
type ArrayValueWriter struct {
	ValueWriter
	// ShowIndexes shows the indexes.
	// Default: false.
	ShowIndexes bool
	// MaxLen is the maximum length of the array.
	// Default: 0 (no limit).
	MaxLen int
}

// NewArrayValueWriter creates a new [ArrayValueWriter] with default values.
func NewArrayValueWriter(vw ValueWriter) *ArrayValueWriter {
	return &ArrayValueWriter{
		ValueWriter: vw,
		ShowIndexes: false,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ArrayValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Array {
		return false
	}
	writeArray(st, v, vw.ShowIndexes, vw.MaxLen, vw.ValueWriter)
	return true
}

func writeArray(st *State, v reflect.Value, showIndexes bool, maxLen int, vw ValueWriter) {
	l := v.Len()
	truncated := false
	if maxLen > 0 && l > maxLen {
		l = maxLen
		truncated = true
	}
	internal.MustWriteString(st.Writer, "{")
	if v.Len() > 0 {
		internal.MustWriteString(st.Writer, "\n")
		st.IndentLevel++
		for i := range l {
			st.writeIndent()
			if showIndexes {
				internal.MustWrite(strconvio.WriteInt(st.Writer, int64(i), 10))
				internal.MustWriteString(st.Writer, ": ")
			}
			must.Handle(vw(st, v.Index(i)))
			internal.MustWriteString(st.Writer, ",\n")
		}
		if truncated {
			st.writeIndent()
			writeTruncated(st.Writer)
			internal.MustWriteString(st.Writer, "\n")
		}
		st.IndentLevel--
		st.writeIndent()
	}
	internal.MustWriteString(st.Writer, "}")
}
