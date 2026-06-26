package pretty

import (
	"reflect"
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
}

// NewChanWriter creates a new [ChanWriter] with default values.
func NewChanWriter(vw ValueWriter) *ChanWriter {
	return &ChanWriter{
		ValueWriter: vw,
		ShowLen:     true,
		ShowCap:     true,
		ShowAddr:    false,
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
	}.write(st)
	return true
}

// Supports implements [SupportChecker].
func (vw *ChanWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Chan {
		res = vw
	}
	return res
}
