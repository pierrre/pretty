package pretty

import (
	"fmt"
	"io"
	"reflect"
	"runtime"

	"github.com/pierrre/pretty/internal"
)

// PanicRecoverValueWriter is a [ValueWriter] that recovers from panics.
//
// It should be created with [NewPanicRecoverValueWriter].
type PanicRecoverValueWriter struct {
	ValueWriter
	// ShowStack shows the stack trace.
	// Default: true.
	ShowStack bool
}

// NewPanicRecoverValueWriter creates a new [PanicRecoverValueWriter] with default values.
func NewPanicRecoverValueWriter(vw ValueWriter) *PanicRecoverValueWriter {
	return &PanicRecoverValueWriter{
		ValueWriter: vw,
		ShowStack:   true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *PanicRecoverValueWriter) WriteValue(st *State, v reflect.Value) (handled bool) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		handled = true
		_, _ = internal.WriteString(st.Writer, "<panic>: ")
		switch r := r.(type) {
		case string:
			_, _ = internal.WriteString(st.Writer, r)
		case error:
			_, _ = internal.WriteString(st.Writer, r.Error())
		default:
			_, _ = fmt.Fprint(st.Writer, r)
		}
		_, _ = internal.WriteString(st.Writer, "\n")
		if vw.ShowStack {
			writeStack(st.Writer)
		}
	}()
	return vw.ValueWriter(st, v)
}

func writeStack(w io.Writer) {
	bp := bytesPool.Get()
	defer bytesPool.Put(bp)
	for {
		n := runtime.Stack(*bp, false)
		if n < len(*bp) {
			_, _ = w.Write((*bp)[:n])
			return
		}
		*bp = make([]byte, 2*len(*bp))
	}
}
