package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// ArrayWriter is a [ValueWriter] that handles array values.
//
// It should be created with [NewArrayWriter].
type ArrayWriter struct {
	ValueWriter
	// ShowIndexes shows the indexes.
	// Default: false.
	ShowIndexes bool
	// MaxLen is the maximum length of the array.
	// Default: 0 (no limit).
	MaxLen int
}

// NewArrayWriter creates a new [ArrayWriter] with default values.
func NewArrayWriter(vw ValueWriter) *ArrayWriter {
	return &ArrayWriter{
		ValueWriter: vw,
		ShowIndexes: false,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ArrayWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Array {
		return false
	}
	writeArray(st, v, vw.ShowIndexes, vw.MaxLen, vw.ValueWriter)
	return true
}

// Supports implements [SupportChecker].
func (vw *ArrayWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Array {
		res = vw
	}
	return res
}

func writeArray(st *State, v reflect.Value, showIndexes bool, maxLen int, vw ValueWriter) {
	l := v.Len()
	truncated := false
	if maxLen > 0 && l > maxLen {
		l = maxLen
		truncated = true
	}
	write.MustString(st.Writer, "{")
	if l > 0 {
		write.MustString(st.Writer, "\n")
		st.IndentLevel++
		for i := range l {
			st.WriteIndent()
			if showIndexes {
				write.Must(strconvio.WriteInt(st.Writer, int64(i), 10))
				write.MustString(st.Writer, ": ")
			}
			must.Handle(vw.WriteValue(st, v.Index(i)))
			write.MustString(st.Writer, ",\n")
		}
		if truncated {
			st.WriteIndent()
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
		}
		st.IndentLevel--
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
}
