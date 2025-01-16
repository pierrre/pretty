package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/pretty/internal"
)

func checkNil(w io.Writer, v reflect.Value) bool {
	if v.IsNil() {
		writeNil(w)
		return true
	}
	return false
}

func writeNil(w io.Writer) {
	internal.MustWriteString(w, "<nil>")
}
