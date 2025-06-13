package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// BoolWriter is a [ValueWriter] that handles bool values.
//
// It should be created with [NewBoolWriter].
type BoolWriter struct{}

// NewBoolWriter creates a new [BoolWriter].
func NewBoolWriter() *BoolWriter {
	return &BoolWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *BoolWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Bool {
		return false
	}
	write.Must(strconvio.WriteBool(st.Writer, v.Bool()))
	return true
}

// Supports implements [SupportChecker].
func (vw *BoolWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Bool {
		res = vw
	}
	return res
}
