package pretty

import (
	"iter"
	"reflect"

	"github.com/pierrre/pretty/internal"
)

// IterValueWriter is a [ValueWriter] that handles iterators ([iter.Seq] or [iter.Seq2]).
type IterValueWriter struct {
	ValueWriter
	// ShowKeysInfos shows keys infos (for [iter.Seq2]).
	// Default: false.
	ShowKeysInfos bool
	// MaxLen is the maximum length of the iterator.
	// Default: 0 (no limit).
	MaxLen int
}

// NewIterValueWriter creates a new [IterValueWriter] with default values.
func NewIterValueWriter(vw ValueWriter) *IterValueWriter {
	return &IterValueWriter{
		ValueWriter:   vw,
		ShowKeysInfos: false,
		MaxLen:        0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IterValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Func {
		return false
	}
	if v.IsNil() {
		return false
	}
	typ := v.Type()
	switch {
	case typ.CanSeq():
		vw.writeSeq(st, v.Seq())
		return true
	case typ.CanSeq2():
		vw.writeSeq2(st, v.Seq2())
		return true
	}
	return false
}

func (vw *IterValueWriter) writeSeq(st *State, it iter.Seq[reflect.Value]) {
	first := true
	writeString(st.Writer, "{")
	st.IndentLevel++
	i := 0
	for v := range it {
		if first {
			first = false
			writeString(st.Writer, "\n")
		}
		st.writeIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st.Writer)
			writeString(st.Writer, "\n")
			break
		}
		internal.MustHandle(vw.ValueWriter(st, v))
		writeString(st.Writer, ",\n")
		i++
	}
	st.IndentLevel--
	writeString(st.Writer, "}")
}

func (vw *IterValueWriter) writeSeq2(st *State, it iter.Seq2[reflect.Value, reflect.Value]) {
	first := true
	writeString(st.Writer, "{")
	st.IndentLevel++
	for k, v := range it {
		if first {
			first = false
			writeString(st.Writer, "\n")
		}
		showInfos := st.ShowInfos
		st.ShowInfos = vw.ShowKeysInfos
		st.writeIndent()
		internal.MustHandle(vw.ValueWriter(st, k))
		st.ShowInfos = showInfos
		writeString(st.Writer, ": ")
		internal.MustHandle(vw.ValueWriter(st, v))
		writeString(st.Writer, ",\n")
	}
	st.IndentLevel--
	writeString(st.Writer, "}")
}
