package pretty

import (
	"reflect"
	"strconv"
)

// ChanWriter is a [ValueWriter] that handles chan values.
//
// It should be created with [NewChanWriter].
type ChanWriter struct {
	ValueWriter
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// ShowAddr shows the address.
	// Default: false.
	ShowAddr bool
	// ShowElems shows the elements.
	// It reads elements from the channel and puts them back.
	// If the channel is closed, it does not put elements back into the channel.
	// It only works with bidirectional channels.
	// Default: false.
	ShowElems bool
	// ShowIndexes shows the indexes.
	// Default: false.
	ShowIndexes bool
	// MaxLen is the maximum length of the channel.
	// Default: 0 (no limit).
	MaxLen int
}

// NewChanWriter creates a new [ChanWriter] with default values.
func NewChanWriter(vw ValueWriter) *ChanWriter {
	return &ChanWriter{
		ValueWriter: vw,
		ShowLen:     true,
		ShowCap:     true,
		ShowAddr:    false,
		ShowElems:   false,
		ShowIndexes: false,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ChanWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Chan {
		return false
	}
	if checkNil(st, v) {
		return true
	}
	infos{
		showLen:  vw.ShowLen,
		len:      v.Len(),
		showCap:  vw.ShowCap,
		cap:      v.Cap(),
		showAddr: vw.ShowAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(st)
	if v.Type().ChanDir() == reflect.BothDir && vw.ShowElems {
		vw.writeElems(st, v)
	}
	return true
}

func (vw *ChanWriter) writeElems(st *State, v reflect.Value) {
	l := v.Len()
	truncated := false
	if vw.MaxLen > 0 && l > vw.MaxLen {
		l = vw.MaxLen
		truncated = true
	}
	st.Writer.AppendByte('{')
	if l > 0 {
		st.Writer.AppendByte('\n')
		st.IndentLevel++
		for i := range l {
			vw.writeElem(st, v, i)
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

func (vw *ChanWriter) writeElem(st *State, v reflect.Value, i int) {
	st.WriteIndent()
	if vw.ShowIndexes {
		st.Writer = strconv.AppendInt(st.Writer, int64(i), 10)
		st.Writer.AppendString(": ")
	}
	e, _ := v.Recv()
	vw.ValueWriter.WriteValue(st, e)
	func() {
		defer func() {
			_ = recover()
		}()
		v.Send(e)
	}()
	st.Writer.AppendString(",\n")
}

// Supports implements [SupportChecker].
func (vw *ChanWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Chan {
		res = vw
	}
	return res
}
