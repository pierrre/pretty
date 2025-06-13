package pretty

import (
	"reflect"
	"runtime"

	"github.com/pierrre/go-libs/syncutil"
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
	write.MustString(st.Writer, getFuncName(p))
	return true
}

var funcNameCache syncutil.Map[uintptr, string]

func getFuncName(p uintptr) string {
	name, ok := funcNameCache.Load(p)
	if !ok {
		name = runtime.FuncForPC(p).Name()
		funcNameCache.Store(p, name)
	}
	return name
}

// Supports implements [SupportChecker].
func (vw *FuncWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Func {
		res = vw
	}
	return res
}
