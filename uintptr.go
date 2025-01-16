package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal"
)

// UintptrValueWriter is a [ValueWriter] that handles uintptr values.
//
// It should be created with [NewUintptrValueWriter].
type UintptrValueWriter struct{}

// NewUintptrValueWriter creates a new [UintptrValueWriter].
func NewUintptrValueWriter() *UintptrValueWriter {
	return &UintptrValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *UintptrValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Uintptr {
		return false
	}
	writeUintptr(st.Writer, uintptr(v.Uint()))
	return true
}

func writeUintptr(w io.Writer, p uintptr) {
	internal.MustWriteString(w, "0x")
	internal.MustWrite(strconvio.WriteUint(w, uint64(p), 16))
}
