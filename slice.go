package pretty

import (
	"reflect"
)

// SliceWriter is a [ValueWriter] that handles slice values.
//
// It should be created with [NewSliceWriter].
type SliceWriter struct {
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
	// ShowIndexes shows the indexes.
	// Default: false.
	ShowIndexes bool
	// MaxLen is the maximum length of the slice.
	// Default: 0 (no limit).
	MaxLen int
}

// NewSliceWriter creates a new [SliceWriter] with default values.
func NewSliceWriter(vw ValueWriter) *SliceWriter {
	return &SliceWriter{
		ValueWriter: vw,
		ShowLen:     true,
		ShowCap:     true,
		ShowAddr:    false,
		ShowIndexes: false,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *SliceWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Slice {
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
	writeArray(st, v, vw.ShowIndexes, vw.MaxLen, vw.ValueWriter)
	return true
}

// Supports implements [SupportChecker].
func (vw *SliceWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Slice {
		res = vw
	}
	return res
}
