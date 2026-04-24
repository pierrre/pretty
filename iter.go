package pretty

import (
	"reflect"
)

// IterWriter is a [ValueWriter] that handles [iter.Seq] and [iter.Seq2].
//
// It should be created with [NewIterWriter].
type IterWriter struct {
	Seq  *IterSeqWriter
	Seq2 *IterSeq2Writer
}

// NewIterWriter creates a new [IterWriter] with default values.
func NewIterWriter(vw ValueWriter) *IterWriter {
	return &IterWriter{
		Seq:  NewIterSeqWriter(vw),
		Seq2: NewIterSeq2Writer(vw),
	}
}

// WriteValue implements [ValueWriter].
func (vw *IterWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.Seq != nil && vw.Seq.WriteValue(st, v) {
		return true
	}
	if vw.Seq2 != nil && vw.Seq2.WriteValue(st, v) {
		return true
	}
	return false
}

// Supports implements [SupportChecker].
func (vw *IterWriter) Supports(typ reflect.Type) ValueWriter {
	if w := callSupportCheckerPointer(vw.Seq, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Seq2, typ); w != nil {
		return w
	}
	return nil
}

// IterSeqWriter is a [ValueWriter] that handles [iter.Seq].
//
// It should be created with [NewIterSeqWriter].
type IterSeqWriter struct {
	ValueWriter
	// MaxLen is the maximum length of the iterator.
	// Default: 0 (no limit).
	MaxLen int
}

// NewIterSeqWriter creates a new [IterSeqWriter] with default values.
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
	if checkNil(st, v) {
		return true
	}
	st.Writer.AppendByte('{')
	st.IndentLevel++
	i := 0
	v.Seq()(func(v reflect.Value) bool {
		if i == 0 {
			st.Writer.AppendByte('\n')
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st)
			st.Writer.AppendByte('\n')
			return false
		}
		vw.ValueWriter.WriteValue(st, v)
		st.Writer.AppendString(",\n")
		i++
		return true
	})
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	st.Writer.AppendByte('}')
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
	// ShowKeysInfos shows key information.
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
	if checkNil(st, v) {
		return true
	}
	st.Writer.AppendByte('{')
	st.IndentLevel++
	i := 0
	v.Seq2()(func(k, v reflect.Value) bool {
		if i == 0 {
			st.Writer.AppendByte('\n')
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st)
			st.Writer.AppendByte('\n')
			return false
		}
		showInfos := st.ShowInfos
		st.ShowInfos = vw.ShowKeysInfos
		vw.ValueWriter.WriteValue(st, k)
		st.ShowInfos = showInfos
		st.Writer.AppendString(": ")
		vw.ValueWriter.WriteValue(st, v)
		st.Writer.AppendString(",\n")
		i++
		return true
	})
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	st.Writer.AppendByte('}')
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
