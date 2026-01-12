package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// UintWriter is a [ValueWriter] that handles uint values.
//
// It should be created with [NewUintWriter].
type UintWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewUintWriter creates a new [UintWriter] with default values.
func NewUintWriter() *UintWriter {
	return &UintWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *UintWriter) WriteValue(st *State, v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		return false
	}
	write.Must(strconvio.WriteUint(st.Writer, v.Uint(), vw.Base))
	return true
}

// Supports implements [SupportChecker].
func (vw *UintWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	switch typ.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res = vw
	}
	return res
}
