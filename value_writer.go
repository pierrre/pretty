package pretty

import (
	"io"
	"reflect"
)

// ValueWriter is a function that writes a [reflect.Value] to a [io.Writer].
//
// It returns true if it handles the value, false otherwise.
// If it returns false, it must not write anything.
//
// Implementations must check [reflect.Value.CanInterface] before using [reflect.Value.Interface].
//
// Implentations can assume that the value is valid.
type ValueWriter func(w io.Writer, st State, v reflect.Value) bool

// ValueWriters is a list of [ValueWriter].
//
// They are tried in order until one handles the value.
type ValueWriters []ValueWriter

// WriteValue implements [ValueWriter].
func (vws ValueWriters) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	for _, vw := range vws {
		ok := vw(w, st, v)
		if ok {
			return true
		}
	}
	return false
}

func mustHandle(h bool) {
	if !h {
		panic("not handled")
	}
}
