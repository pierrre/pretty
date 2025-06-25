package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// IterSeqWriter is a [ValueWriter] that handles [iter.Seq].
//
// It should be created with [NewIterSeqWriter].
type IterSeqWriter struct {
	ValueWriter
	// MaxLen is the maximum length of the iterator.
	// Default: 0 (no limit).
	MaxLen int
}

// NewIterSeqWriter creates a new [IterSeq2Writer] with default values.
func NewIterSeqWriter(vw ValueWriter) *IterSeqWriter {
	return &IterSeqWriter{
		ValueWriter: vw,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IterSeqWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Func {
		return false
	}
	if !v.Type().CanSeq() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	write.MustString(st.Writer, "{")
	st.IndentLevel++
	i := 0
	v.Seq()(func(v reflect.Value) bool {
		if i == 0 {
			write.MustString(st.Writer, "\n")
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
			return false
		}
		must.Handle(vw.ValueWriter.WriteValue(st, v))
		write.MustString(st.Writer, ",\n")
		i++
		return true
	})
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

// Supports implements [SupportChecker].
func (vw *IterSeqWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Func && typ.CanSeq() {
		res = vw
	}
	return res
}

// IterSeq2Writer is a [ValueWriter] that handles [iter.Seq2].
//
// It should be created with [NewIterSeq2Writer].
type IterSeq2Writer struct {
	ValueWriter
	// ShowKeysInfos shows keys infos.
	// Default: false.
	ShowKeysInfos bool
	// MaxLen is the maximum length of the iterator.
	// Default: 0 (no limit).
	MaxLen int
}

// NewIterSeq2Writer creates a new [IterSeq2Writer] with default values.
func NewIterSeq2Writer(vw ValueWriter) *IterSeq2Writer {
	return &IterSeq2Writer{
		ValueWriter:   vw,
		ShowKeysInfos: false,
		MaxLen:        0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IterSeq2Writer) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Func {
		return false
	}
	if !v.Type().CanSeq2() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	write.MustString(st.Writer, "{")
	st.IndentLevel++
	i := 0
	v.Seq2()(func(k, v reflect.Value) bool {
		if i == 0 {
			write.MustString(st.Writer, "\n")
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
			return false
		}
		showInfos := st.ShowInfos
		st.ShowInfos = vw.ShowKeysInfos
		must.Handle(vw.ValueWriter.WriteValue(st, k))
		st.ShowInfos = showInfos
		write.MustString(st.Writer, ": ")
		must.Handle(vw.ValueWriter.WriteValue(st, v))
		write.MustString(st.Writer, ",\n")
		i++
		return true
	})
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

// Supports implements [SupportChecker].
func (vw *IterSeq2Writer) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Func && typ.CanSeq2() {
		res = vw
	}
	return res
}
