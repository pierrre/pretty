package pretty

import (
	"io"
	"reflect"
)

// SliceValueWriter is a [ValueWriter] that handles slice values.
//
// It should be created with [NewSliceValueWriter].
type SliceValueWriter struct {
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
	// MaxLen is the maximum length of the slice.
	// Default: 0 (no limit).
	MaxLen int
}

// NewSliceValueWriter creates a new [SliceValueWriter] with default values.
func NewSliceValueWriter(vw ValueWriter) *SliceValueWriter {
	return &SliceValueWriter{
		ValueWriter: vw,
		ShowLen:     true,
		ShowCap:     true,
		ShowAddr:    false,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *SliceValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Kind() != reflect.Slice {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	infos{
		showLen:  vw.ShowLen,
		len:      v.Len(),
		showCap:  vw.ShowCap,
		cap:      v.Cap(),
		showAddr: vw.ShowAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(w)
	writeArray(w, st, v, vw.MaxLen, vw.ValueWriter)
	return true
}
