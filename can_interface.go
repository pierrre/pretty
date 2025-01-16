package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// CanInterfaceValueWriter is a [ValueWriter] that attempts to convert the [reflect.Value] so it can be used with [reflect.Value.Interface].
//
// It should be created with [NewCanInterfaceValueWriter].
type CanInterfaceValueWriter struct {
	ValueWriter
}

func NewCanInterfaceValueWriter(vw ValueWriter) *CanInterfaceValueWriter {
	return &CanInterfaceValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *CanInterfaceValueWriter) WriteValue(st *State, v reflect.Value) bool {
	v, _ = reflectutil.ConvertValueCanInterface(v)
	return vw.ValueWriter.WriteValue(st, v)
}
