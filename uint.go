package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
)

// UintValueWriter is a [ValueWriter] that handles uint values.
//
// It should be created with [NewUintValueWriter].
type UintValueWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewUintValueWriter creates a new [UintValueWriter] with default values.
func NewUintValueWriter() *UintValueWriter {
	return &UintValueWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *UintValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles uint.
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mustWrite(strconvio.WriteUint(w, v.Uint(), vw.Base))
		return true
	}
	return false
}