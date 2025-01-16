package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/pretty/internal/write"
)

func checkNil(w io.Writer, v reflect.Value) bool {
	if v.IsNil() {
		writeNil(w)
		return true
	}
	return false
}

func writeNil(w io.Writer) {
	write.MustString(w, "<nil>")
}
