package pretty

import (
	"fmt"
	"io"
	"reflect"
	"runtime/debug"
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
func (vw *PanicRecoverValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) (handled bool) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		handled = true
		_, _ = writeStringErr(w, "<panic>: ")
		switch r := r.(type) {
		case string:
			_, _ = writeStringErr(w, r)
		case error:
			_, _ = writeStringErr(w, r.Error())
		default:
			_, _ = fmt.Fprint(w, r)
		}
		_, _ = writeStringErr(w, "\n")
		if vw.ShowStack {
			_, _ = w.Write(debug.Stack())
		}
	}()
	return vw.ValueWriter(w, st, v)
}