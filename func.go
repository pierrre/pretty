package pretty

import (
	"reflect"
	"runtime"

	"github.com/pierrre/pretty/internal/write"
)

// FuncWriter is a [ValueWriter] that handles function values.
//
// It should be created with [NewFuncWriter].
type FuncWriter struct {
	// ShowAddr shows the address.
	// Default: false.
	ShowAddr bool
}

// NewFuncWriter creates a new [FuncWriter] with default values.
func NewFuncWriter() *FuncWriter {
	return &FuncWriter{
		ShowAddr: false,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FuncWriter) WriteValue(st *State, v reflect.Value) bool {
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
	write.MustString(st.Writer, name)
	return true
}

// Supports implements [SupportChecker].
func (vw *FuncWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Func {
		res = vw
	}
	return res
}
