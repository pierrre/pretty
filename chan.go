package pretty

import (
	"reflect"
)

// ChanValueWriter is a [ValueWriter] that handles chan values.
//
// It should be created with [NewChanValueWriter].
type ChanValueWriter struct {
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

// NewChanValueWriter creates a new [ChanValueWriter] with default values.
func NewChanValueWriter() *ChanValueWriter {
	return &ChanValueWriter{
		ShowLen:  true,
		ShowCap:  true,
		ShowAddr: false,
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
	}.write(st)
	return true
}
