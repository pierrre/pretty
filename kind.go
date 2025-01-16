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
	BaseInvalid       *InvalidValueWriter
	BaseBool          *BoolValueWriter
	BaseInt           *IntValueWriter
	BaseUint          *UintValueWriter
	BaseUintptr       *UintptrValueWriter
	BaseFloat         *FloatValueWriter
	BaseComplex       *ComplexValueWriter
	BaseArray         *ArrayValueWriter
	BaseChan          *ChanValueWriter
	BaseFunc          *FuncValueWriter
	BaseInterface     *InterfaceValueWriter
	BaseMap           *MapValueWriter
	BasePointer       *PointerValueWriter
	BaseSlice         *SliceValueWriter
	BaseString        *StringValueWriter
	BaseStruct        *StructValueWriter
	BaseUnsafePointer *UnsafePointerValueWriter

	ValueWriters [kindsCount]ValueWriter
}

// NewKindValueWriter creates a new [KindValueWriter] with default values.
func NewKindValueWriter(vw ValueWriter) *KindValueWriter {
	kindVW := &KindValueWriter{
		BaseInvalid:       NewInvalidValueWriter(),
		BaseBool:          NewBoolValueWriter(),
		BaseInt:           NewIntValueWriter(),
		BaseUint:          NewUintValueWriter(),
		BaseUintptr:       NewUintptrValueWriter(),
		BaseFloat:         NewFloatValueWriter(),
		BaseComplex:       NewComplexValueWriter(),
		BaseArray:         NewArrayValueWriter(vw),
		BaseChan:          NewChanValueWriter(vw),
		BaseFunc:          NewFuncValueWriter(),
		BaseInterface:     NewInterfaceValueWriter(vw),
		BaseMap:           NewMapValueWriter(vw),
		BasePointer:       NewPointerValueWriter(vw),
		BaseSlice:         NewSliceValueWriter(vw),
		BaseString:        NewStringValueWriter(),
		BaseStruct:        NewStructValueWriter(vw),
		BaseUnsafePointer: NewUnsafePointerValueWriter(),
	}
	kindVW.ValueWriters = [kindsCount]ValueWriter{
		reflect.Invalid:       ValueWriterFunc(kindVW.writeInvalid),
		reflect.Bool:          ValueWriterFunc(kindVW.writeBool),
		reflect.Int:           ValueWriterFunc(kindVW.writeInt),
		reflect.Int8:          ValueWriterFunc(kindVW.writeInt),
		reflect.Int16:         ValueWriterFunc(kindVW.writeInt),
		reflect.Int32:         ValueWriterFunc(kindVW.writeInt),
		reflect.Int64:         ValueWriterFunc(kindVW.writeInt),
		reflect.Uint:          ValueWriterFunc(kindVW.writeUint),
		reflect.Uint8:         ValueWriterFunc(kindVW.writeUint),
		reflect.Uint16:        ValueWriterFunc(kindVW.writeUint),
		reflect.Uint32:        ValueWriterFunc(kindVW.writeUint),
		reflect.Uint64:        ValueWriterFunc(kindVW.writeUint),
		reflect.Uintptr:       ValueWriterFunc(kindVW.writeUintptr),
		reflect.Float32:       ValueWriterFunc(kindVW.writeFloat),
		reflect.Float64:       ValueWriterFunc(kindVW.writeFloat),
		reflect.Complex64:     ValueWriterFunc(kindVW.writeComplex),
		reflect.Complex128:    ValueWriterFunc(kindVW.writeComplex),
		reflect.Array:         ValueWriterFunc(kindVW.writeArray),
		reflect.Chan:          ValueWriterFunc(kindVW.writeChan),
		reflect.Func:          ValueWriterFunc(kindVW.writeFunc),
		reflect.Interface:     ValueWriterFunc(kindVW.writeInterface),
		reflect.Map:           ValueWriterFunc(kindVW.writeMap),
		reflect.Pointer:       ValueWriterFunc(kindVW.writePointer),
		reflect.Slice:         ValueWriterFunc(kindVW.writeSlice),
		reflect.String:        ValueWriterFunc(kindVW.writeString),
		reflect.Struct:        ValueWriterFunc(kindVW.writeStruct),
		reflect.UnsafePointer: ValueWriterFunc(kindVW.writeUnsafePointer),
	}
	return kindVW
}

// WriteValue implements [ValueWriter].
func (vw *KindValueWriter) WriteValue(st *State, v reflect.Value) bool {
	return vw.ValueWriters[v.Kind()].WriteValue(st, v)
}

func (vw *KindValueWriter) writeInvalid(st *State, v reflect.Value) bool {
	return vw.BaseInvalid.WriteValue(st, v)
}

func (vw *KindValueWriter) writeBool(st *State, v reflect.Value) bool {
	return vw.BaseBool.WriteValue(st, v)
}

func (vw *KindValueWriter) writeInt(st *State, v reflect.Value) bool {
	return vw.BaseInt.WriteValue(st, v)
}

func (vw *KindValueWriter) writeUint(st *State, v reflect.Value) bool {
	return vw.BaseUint.WriteValue(st, v)
}

func (vw *KindValueWriter) writeUintptr(st *State, v reflect.Value) bool {
	return vw.BaseUintptr.WriteValue(st, v)
}

func (vw *KindValueWriter) writeFloat(st *State, v reflect.Value) bool {
	return vw.BaseFloat.WriteValue(st, v)
}

func (vw *KindValueWriter) writeComplex(st *State, v reflect.Value) bool {
	return vw.BaseComplex.WriteValue(st, v)
}

func (vw *KindValueWriter) writeArray(st *State, v reflect.Value) bool {
	return vw.BaseArray.WriteValue(st, v)
}

func (vw *KindValueWriter) writeChan(st *State, v reflect.Value) bool {
	return vw.BaseChan.WriteValue(st, v)
}

func (vw *KindValueWriter) writeFunc(st *State, v reflect.Value) bool {
	return vw.BaseFunc.WriteValue(st, v)
}

func (vw *KindValueWriter) writeInterface(st *State, v reflect.Value) bool {
	return vw.BaseInterface.WriteValue(st, v)
}

func (vw *KindValueWriter) writeMap(st *State, v reflect.Value) bool {
	return vw.BaseMap.WriteValue(st, v)
}

func (vw *KindValueWriter) writePointer(st *State, v reflect.Value) bool {
	return vw.BasePointer.WriteValue(st, v)
}

func (vw *KindValueWriter) writeSlice(st *State, v reflect.Value) bool {
	return vw.BaseSlice.WriteValue(st, v)
}

func (vw *KindValueWriter) writeString(st *State, v reflect.Value) bool {
	return vw.BaseString.WriteValue(st, v)
}

func (vw *KindValueWriter) writeStruct(st *State, v reflect.Value) bool {
	return vw.BaseStruct.WriteValue(st, v)
}

func (vw *KindValueWriter) writeUnsafePointer(st *State, v reflect.Value) bool {
	return vw.BaseUnsafePointer.WriteValue(st, v)
}
