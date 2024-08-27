package pretty

import (
	"io"
	"reflect"
	"runtime"
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
func (vw *FuncValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() != reflect.Func {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	p := uintptr(v.UnsafePointer())
	infos{
		showAddr: vw.ShowAddr,
		addr:     p,
	}.writeWithTrailingSpace(w)
	name := runtime.FuncForPC(p).Name()
	writeString(w, name)
	return true
}
