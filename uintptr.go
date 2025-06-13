package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// UintptrWriter is a [ValueWriter] that handles uintptr values.
//
// It should be created with [NewUintptrWriter].
type UintptrWriter struct{}

// NewUintptrWriter creates a new [UintptrWriter].
func NewUintptrWriter() *UintptrWriter {
	return &UintptrWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *UintptrWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Uintptr {
		return false
	}
	writeUintptr(st.Writer, uintptr(v.Uint()))
	return true
}

// Supports implements [SupportChecker].
func (vw *UintptrWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Uintptr {
		res = vw
	}
	return res
}

func writeUintptr(w io.Writer, p uintptr) {
	write.MustString(w, "0x")
	write.Must(strconvio.WriteUint(w, uint64(p), 16))
}
