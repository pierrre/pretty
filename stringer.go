package pretty

import (
	"fmt"
	"reflect"
)

var stringerType = reflect.TypeFor[fmt.Stringer]()

// StringerValueWriter is a [ValueWriter] that handles [fmt.Stringer].
//
// It should be created with [NewStringerValueWriter].
type StringerValueWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// Quote quotes the string.
	// Default: true.
	Quote bool
	// MaxLen is the maximum length of the string.
	// Default: 0 (no limit).
	MaxLen int
}

// NewStringerValueWriter creates a new [StringerValueWriter].
func NewStringerValueWriter() *StringerValueWriter {
	return &StringerValueWriter{
		ShowLen: true,
		Quote:   true,
		MaxLen:  0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StringerValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !v.Type().Implements(stringerType) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if v.Type() == reflectValueType {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	sr := v.Interface().(fmt.Stringer) //nolint:forcetypeassert // Checked above.
	s := sr.String()
	writeArrowWrappedString(st.Writer, ".String() ")
	writeStringValue(st, s, vw.ShowLen, false, 0, vw.Quote, vw.MaxLen)
	return true
}
