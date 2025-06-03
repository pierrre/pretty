package pretty

import (
	"reflect"
)

const kindsCount = reflect.UnsafePointer + 1

// KindValueWriter is a [ValueWriter] that writes the value with the kind-specific [ValueWriter].
//
// It should be created with [NewKindValueWriter].
type KindValueWriter struct {
	// The following fields are the default [ValueWriter] for each group of kinds.
	Invalid       InvalidValueWriter
	Bool          BoolValueWriter
	Int           IntValueWriter
	Uint          UintValueWriter
	Uintptr       UintptrValueWriter
	Float         FloatValueWriter
	Complex       ComplexValueWriter
	Array         ArrayValueWriter
	Chan          ChanValueWriter
	Func          FuncValueWriter
	Interface     InterfaceValueWriter
	Map           MapValueWriter
	Pointer       PointerValueWriter
	Slice         SliceValueWriter
	String        StringValueWriter
	Struct        StructValueWriter
	UnsafePointer UnsafePointerValueWriter

	// ValueWriters is the list of [ValueWriter] indexed by [reflect.Kind].
	ValueWriters [kindsCount]ValueWriter
}

// NewKindValueWriter creates a new [KindValueWriter] with default values.
func NewKindValueWriter(vw ValueWriter) *KindValueWriter {
	kindVW := &KindValueWriter{
		Invalid:       *NewInvalidValueWriter(),
		Bool:          *NewBoolValueWriter(),
		Int:           *NewIntValueWriter(),
		Uint:          *NewUintValueWriter(),
		Uintptr:       *NewUintptrValueWriter(),
		Float:         *NewFloatValueWriter(),
		Complex:       *NewComplexValueWriter(),
		Array:         *NewArrayValueWriter(vw),
		Chan:          *NewChanValueWriter(vw),
		Func:          *NewFuncValueWriter(),
		Interface:     *NewInterfaceValueWriter(vw),
		Map:           *NewMapValueWriter(vw),
		Pointer:       *NewPointerValueWriter(vw),
		Slice:         *NewSliceValueWriter(vw),
		String:        *NewStringValueWriter(),
		Struct:        *NewStructValueWriter(vw),
		UnsafePointer: *NewUnsafePointerValueWriter(),
	}
	kindVW.ValueWriters = [kindsCount]ValueWriter{
		reflect.Invalid:       &kindVW.Invalid,
		reflect.Bool:          &kindVW.Bool,
		reflect.Int:           &kindVW.Int,
		reflect.Int8:          &kindVW.Int,
		reflect.Int16:         &kindVW.Int,
		reflect.Int32:         &kindVW.Int,
		reflect.Int64:         &kindVW.Int,
		reflect.Uint:          &kindVW.Uint,
		reflect.Uint8:         &kindVW.Uint,
		reflect.Uint16:        &kindVW.Uint,
		reflect.Uint32:        &kindVW.Uint,
		reflect.Uint64:        &kindVW.Uint,
		reflect.Uintptr:       &kindVW.Uintptr,
		reflect.Float32:       &kindVW.Float,
		reflect.Float64:       &kindVW.Float,
		reflect.Complex64:     &kindVW.Complex,
		reflect.Complex128:    &kindVW.Complex,
		reflect.Array:         &kindVW.Array,
		reflect.Chan:          &kindVW.Chan,
		reflect.Func:          &kindVW.Func,
		reflect.Interface:     &kindVW.Interface,
		reflect.Map:           &kindVW.Map,
		reflect.Pointer:       &kindVW.Pointer,
		reflect.Slice:         &kindVW.Slice,
		reflect.String:        &kindVW.String,
		reflect.Struct:        &kindVW.Struct,
		reflect.UnsafePointer: &kindVW.UnsafePointer,
	}
	return kindVW
}

// WriteValue implements [ValueWriter].
func (vw *KindValueWriter) WriteValue(st *State, v reflect.Value) bool {
	return vw.ValueWriters[v.Kind()].WriteValue(st, v)
}
