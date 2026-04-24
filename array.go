package pretty

import (
	"reflect"
	"strconv"
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
	st.Writer.AppendByte('{')
	if l > 0 {
		st.Writer.AppendByte('\n')
		st.IndentLevel++
		for i := range l {
			st.WriteIndent()
			if showIndexes {
				st.Writer = strconv.AppendInt(st.Writer, int64(i), 10)
				st.Writer.AppendString(": ")
			}
			vw.WriteValue(st, v.Index(i))
			st.Writer.AppendString(",\n")
		}
		if truncated {
			st.WriteIndent()
			writeTruncated(st)
			st.Writer.AppendByte('\n')
		}
		st.IndentLevel--
		st.WriteIndent()
	}
	st.Writer.AppendByte('}')
}
