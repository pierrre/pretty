package pretty

import (
	"reflect"
)

// ValueWriter writes a [reflect.Value] to a [io.Writer].
//
// It returns true if it handles the value, false otherwise.
// If it returns false, it must not write anything.
//
// Implementations must check [reflect.Value.CanInterface] before using [reflect.Value.Interface].
//
// Implentations can assume that the value is valid.
type ValueWriter interface {
	WriteValue(st *State, v reflect.Value) bool
}

// ValueWriterFunc is a [ValueWriter] function.
type ValueWriterFunc func(st *State, v reflect.Value) bool

// WriteValue implements [ValueWriter].
func (f ValueWriterFunc) WriteValue(st *State, v reflect.Value) bool {
	return f(st, v)
}

// ValueWriters is a list of [ValueWriter].
//
// They are tried in order until one handles the value.
type ValueWriters []ValueWriter

// WriteValue implements [ValueWriter].
func (vws ValueWriters) WriteValue(st *State, v reflect.Value) bool {
	for _, vw := range vws {
		ok := vw.WriteValue(st, v)
		if ok {
			return true
		}
	}
	return false
}

// Supports implements [SupportChecker].
func (vws ValueWriters) Supports(typ reflect.Type) ValueWriter {
	for _, vw := range vws {
		f := supportsValueWriter(typ, vw)
		if f != nil {
			return f
		}
	}
	return nil
}
