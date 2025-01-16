package pretty

import (
	"reflect"
	"runtime"

	"github.com/pierrre/pretty/internal"
)

// FuncValueWriter is a [ValueWriter] that handles function values.
//
// It should be created with [NewFuncValueWriter].
type FuncValueWriter struct {
	// ShowAddr shows the address.
	// Default: false.
	ShowAddr bool
}

// NewFuncValueWriter creates a new [FuncValueWriter] with default values.
func NewFuncValueWriter() *FuncValueWriter {
	return &FuncValueWriter{
		ShowAddr: false,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FuncValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Func {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	p := uintptr(v.UnsafePointer())
	infos{
		showAddr: vw.ShowAddr,
		addr:     p,
	}.writeWithTrailingSpace(st)
	name := runtime.FuncForPC(p).Name()
	internal.MustWriteString(st.Writer, name)
	return true
}
