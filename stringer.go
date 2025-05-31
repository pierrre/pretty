package pretty

import (
	"fmt"
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/itfassert"
)

var stringerImplementsCache = reflectutil.NewImplementsCacheFor[fmt.Stringer]()

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
	typ := v.Type()
	if !stringerImplementsCache.ImplementedBy(typ) {
		return false
	}
	if typ == reflectValueType {
		return false
	}
	sr, ok := itfassert.Assert[fmt.Stringer](v)
	if !ok {
		return false
	}
	s := sr.String()
	writeArrowWrappedString(st.Writer, ".String() ")
	writeStringValue(st, s, vw.ShowLen, false, 0, vw.Quote, vw.MaxLen)
	return true
}
