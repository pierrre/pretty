package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// ChanValueWriter is a [ValueWriter] that handles chan values.
//
// It should be created with [NewChanValueWriter].
type ChanValueWriter struct {
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
	// It reads the elements from the channel and put them back.
	// If the channel is closed, it doesn't put the elements back to the channel.
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

// NewChanValueWriter creates a new [ChanValueWriter] with default values.
func NewChanValueWriter(vw ValueWriter) *ChanValueWriter {
	return &ChanValueWriter{
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
func (vw *ChanValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Chan {
		return false
	}
	if checkNil(st.Writer, v) {
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

func (vw *ChanValueWriter) writeElems(st *State, v reflect.Value) {
	l := v.Len()
	truncated := false
	if vw.MaxLen > 0 && l > vw.MaxLen {
		l = vw.MaxLen
		truncated = true
	}
	write.MustString(st.Writer, "{")
	if v.Len() > 0 {
		write.MustString(st.Writer, "\n")
		st.IndentLevel++
		for i := range l {
			vw.writeElem(st, v, i)
		}
		if truncated {
			st.writeIndent()
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
		}
		st.IndentLevel--
		st.writeIndent()
	}
	write.MustString(st.Writer, "}")
}

func (vw *ChanValueWriter) writeElem(st *State, v reflect.Value, i int) {
	st.writeIndent()
	if vw.ShowIndexes {
		write.Must(strconvio.WriteInt(st.Writer, int64(i), 10))
		write.MustString(st.Writer, ": ")
	}
	e, _ := v.Recv()
	must.Handle(vw.ValueWriter.WriteValue(st, e))
	func() {
		defer func() {
			_ = recover()
		}()
		v.Send(e)
	}()
	write.MustString(st.Writer, ",\n")
}
