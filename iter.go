package pretty

import (
	"iter"
	"reflect"

	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// IterWriter is a [ValueWriter] that handles iterators ([iter.Seq] or [iter.Seq2]).
type IterWriter struct {
	ValueWriter
	// ShowKeysInfos shows keys infos (for [iter.Seq2]).
	// Default: false.
	ShowKeysInfos bool
	// MaxLen is the maximum length of the iterator.
	// Default: 0 (no limit).
	MaxLen int
}

// NewIterWriter creates a new [IterWriter] with default values.
func NewIterWriter(vw ValueWriter) *IterWriter {
	return &IterWriter{
		ValueWriter:   vw,
		ShowKeysInfos: false,
		MaxLen:        0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IterWriter) WriteValue(st *State, v reflect.Value) bool {
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

func (vw *IterWriter) writeSeq(st *State, it iter.Seq[reflect.Value]) {
	write.MustString(st.Writer, "{")
	st.IndentLevel++
	i := 0
	for v := range it {
		if i == 0 {
			write.MustString(st.Writer, "\n")
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
			break
		}
		must.Handle(vw.ValueWriter.WriteValue(st, v))
		write.MustString(st.Writer, ",\n")
		i++
	}
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
}

func (vw *IterWriter) writeSeq2(st *State, it iter.Seq2[reflect.Value, reflect.Value]) {
	write.MustString(st.Writer, "{")
	st.IndentLevel++
	i := 0
	for k, v := range it {
		if i == 0 {
			write.MustString(st.Writer, "\n")
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
			break
		}
		showInfos := st.ShowInfos
		st.ShowInfos = vw.ShowKeysInfos
		must.Handle(vw.ValueWriter.WriteValue(st, k))
		st.ShowInfos = showInfos
		write.MustString(st.Writer, ": ")
		must.Handle(vw.ValueWriter.WriteValue(st, v))
		write.MustString(st.Writer, ",\n")
		i++
	}
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
}

// Supports implements [SupportChecker].
func (vw *IterWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Func && (typ.CanSeq() || typ.CanSeq2()) {
		res = vw
	}
	return res
}
