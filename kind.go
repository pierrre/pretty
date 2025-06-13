package pretty

import (
	"reflect"
)

const kindsCount = reflect.UnsafePointer + 1

// KindWriter is a [ValueWriter] that writes the value with the kind-specific [ValueWriter].
//
// It should be created with [NewKindWriter].
type KindWriter struct {
	// The following fields are the default [ValueWriter] for each group of kinds.
	Invalid       *InvalidWriter
	Bool          *BoolWriter
	Int           *IntWriter
	Uint          *UintWriter
	Uintptr       *UintptrWriter
	Float         *FloatWriter
	Complex       *ComplexWriter
	Array         *ArrayWriter
	Chan          *ChanWriter
	Func          *FuncWriter
	Interface     *InterfaceWriter
	Map           *MapWriter
	Pointer       *PointerWriter
	Slice         *SliceWriter
	String        *StringWriter
	Struct        *StructWriter
	UnsafePointer *UnsafePointerWriter

	// ValueWriters is the list of [ValueWriter] indexed by [reflect.Kind].
	ValueWriters [kindsCount]ValueWriter
}

// NewKindWriter creates a new [KindWriter] with default values.
func NewKindWriter(vw ValueWriter) *KindWriter {
	kindVW := &KindWriter{
		Invalid:       NewInvalidWriter(),
		Bool:          NewBoolWriter(),
		Int:           NewIntWriter(),
		Uint:          NewUintWriter(),
		Uintptr:       NewUintptrWriter(),
		Float:         NewFloatWriter(),
		Complex:       NewComplexWriter(),
		Array:         NewArrayWriter(vw),
		Chan:          NewChanWriter(vw),
		Func:          NewFuncWriter(),
		Interface:     NewInterfaceWriter(vw),
		Map:           NewMapWriter(vw),
		Pointer:       NewPointerWriter(vw),
		Slice:         NewSliceWriter(vw),
		String:        NewStringWriter(),
		Struct:        NewStructWriter(vw),
		UnsafePointer: NewUnsafePointerWriter(),
	}
	kindVW.ValueWriters = [kindsCount]ValueWriter{
		reflect.Invalid:       kindVW.Invalid,
		reflect.Bool:          kindVW.Bool,
		reflect.Int:           kindVW.Int,
		reflect.Int8:          kindVW.Int,
		reflect.Int16:         kindVW.Int,
		reflect.Int32:         kindVW.Int,
		reflect.Int64:         kindVW.Int,
		reflect.Uint:          kindVW.Uint,
		reflect.Uint8:         kindVW.Uint,
		reflect.Uint16:        kindVW.Uint,
		reflect.Uint32:        kindVW.Uint,
		reflect.Uint64:        kindVW.Uint,
		reflect.Uintptr:       kindVW.Uintptr,
		reflect.Float32:       kindVW.Float,
		reflect.Float64:       kindVW.Float,
		reflect.Complex64:     kindVW.Complex,
		reflect.Complex128:    kindVW.Complex,
		reflect.Array:         kindVW.Array,
		reflect.Chan:          kindVW.Chan,
		reflect.Func:          kindVW.Func,
		reflect.Interface:     kindVW.Interface,
		reflect.Map:           kindVW.Map,
		reflect.Pointer:       kindVW.Pointer,
		reflect.Slice:         kindVW.Slice,
		reflect.String:        kindVW.String,
		reflect.Struct:        kindVW.Struct,
		reflect.UnsafePointer: kindVW.UnsafePointer,
	}
	return kindVW
}

// WriteValue implements [ValueWriter].
func (vw *KindWriter) WriteValue(st *State, v reflect.Value) bool {
	return vw.ValueWriters[v.Kind()].WriteValue(st, v)
}

// Supports implements [SupportChecker].
func (vw *KindWriter) Supports(typ reflect.Type) ValueWriter {
	return supportsValueWriter(typ, vw.ValueWriters[typ.Kind()])
}
