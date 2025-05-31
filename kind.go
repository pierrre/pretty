package pretty

import (
	"reflect"
)

const kindsCount = reflect.UnsafePointer + 1

// KindValueWriter is a [ValueWriter] that writes the value with the kind-specific [ValueWriter].
//
// The "Base*" fields are the default [ValueWriter] for each group of kinds.
// They can be configured.
//
// The [ValueWriters] fields is the list of [ValueWriter] indexed by [reflect.Kind].
//
// It should be created with [NewKindValueWriter].
type KindValueWriter struct {
	BaseInvalid       InvalidValueWriter
	BaseBool          BoolValueWriter
	BaseInt           IntValueWriter
	BaseUint          UintValueWriter
	BaseUintptr       UintptrValueWriter
	BaseFloat         FloatValueWriter
	BaseComplex       ComplexValueWriter
	BaseArray         ArrayValueWriter
	BaseChan          ChanValueWriter
	BaseFunc          FuncValueWriter
	BaseInterface     InterfaceValueWriter
	BaseMap           MapValueWriter
	BasePointer       PointerValueWriter
	BaseSlice         SliceValueWriter
	BaseString        StringValueWriter
	BaseStruct        StructValueWriter
	BaseUnsafePointer UnsafePointerValueWriter

	ValueWriters [kindsCount]ValueWriter
}

// NewKindValueWriter creates a new [KindValueWriter] with default values.
func NewKindValueWriter(vw ValueWriter) *KindValueWriter {
	kindVW := &KindValueWriter{
		BaseInvalid:       *NewInvalidValueWriter(),
		BaseBool:          *NewBoolValueWriter(),
		BaseInt:           *NewIntValueWriter(),
		BaseUint:          *NewUintValueWriter(),
		BaseUintptr:       *NewUintptrValueWriter(),
		BaseFloat:         *NewFloatValueWriter(),
		BaseComplex:       *NewComplexValueWriter(),
		BaseArray:         *NewArrayValueWriter(vw),
		BaseChan:          *NewChanValueWriter(vw),
		BaseFunc:          *NewFuncValueWriter(),
		BaseInterface:     *NewInterfaceValueWriter(vw),
		BaseMap:           *NewMapValueWriter(vw),
		BasePointer:       *NewPointerValueWriter(vw),
		BaseSlice:         *NewSliceValueWriter(vw),
		BaseString:        *NewStringValueWriter(),
		BaseStruct:        *NewStructValueWriter(vw),
		BaseUnsafePointer: *NewUnsafePointerValueWriter(),
	}
	kindVW.ValueWriters = [kindsCount]ValueWriter{
		reflect.Invalid:       &kindVW.BaseInvalid,
		reflect.Bool:          &kindVW.BaseBool,
		reflect.Int:           &kindVW.BaseInt,
		reflect.Int8:          &kindVW.BaseInt,
		reflect.Int16:         &kindVW.BaseInt,
		reflect.Int32:         &kindVW.BaseInt,
		reflect.Int64:         &kindVW.BaseInt,
		reflect.Uint:          &kindVW.BaseUint,
		reflect.Uint8:         &kindVW.BaseUint,
		reflect.Uint16:        &kindVW.BaseUint,
		reflect.Uint32:        &kindVW.BaseUint,
		reflect.Uint64:        &kindVW.BaseUint,
		reflect.Uintptr:       &kindVW.BaseUintptr,
		reflect.Float32:       &kindVW.BaseFloat,
		reflect.Float64:       &kindVW.BaseFloat,
		reflect.Complex64:     &kindVW.BaseComplex,
		reflect.Complex128:    &kindVW.BaseComplex,
		reflect.Array:         &kindVW.BaseArray,
		reflect.Chan:          &kindVW.BaseChan,
		reflect.Func:          &kindVW.BaseFunc,
		reflect.Interface:     &kindVW.BaseInterface,
		reflect.Map:           &kindVW.BaseMap,
		reflect.Pointer:       &kindVW.BasePointer,
		reflect.Slice:         &kindVW.BaseSlice,
		reflect.String:        &kindVW.BaseString,
		reflect.Struct:        &kindVW.BaseStruct,
		reflect.UnsafePointer: &kindVW.BaseUnsafePointer,
	}
	return kindVW
}

// WriteValue implements [ValueWriter].
func (vw *KindValueWriter) WriteValue(st *State, v reflect.Value) bool {
	return vw.ValueWriters[v.Kind()].WriteValue(st, v)
}
