package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// CanInterfaceWriter is a [ValueWriter] that attempts to convert the [reflect.Value] so it can be used with [reflect.Value.Interface].
//
// It should be created with [NewCanInterfaceWriter].
type CanInterfaceWriter struct {
	ValueWriter
}

// NewCanInterfaceWriter creates a new CanInterfaceValueWriter.
func NewCanInterfaceWriter(vw ValueWriter) *CanInterfaceWriter {
	return &CanInterfaceWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *CanInterfaceWriter) WriteValue(st *State, v reflect.Value) bool {
	v, _ = reflectutil.ConvertValueCanInterface(v)
	return vw.ValueWriter.WriteValue(st, v)
}
