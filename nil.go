package pretty

import (
	"io"
	"reflect"
)

func checkNil(w io.Writer, v reflect.Value) bool {
	if v.IsNil() {
		writeNil(w)
		return true
	}
	return false
}

func writeNil(w io.Writer) {
	writeString(w, "<nil>")
}
