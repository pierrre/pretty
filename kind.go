package pretty

import (
	"io"
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
		BaseChan:          NewChanValueWriter(),
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
		reflect.Invalid:       kindVW.writeInvalid,
		reflect.Bool:          kindVW.writeBool,
		reflect.Int:           kindVW.writeInt,
		reflect.Int8:          kindVW.writeInt,
		reflect.Int16:         kindVW.writeInt,
		reflect.Int32:         kindVW.writeInt,
		reflect.Int64:         kindVW.writeInt,
		reflect.Uint:          kindVW.writeUint,
		reflect.Uint8:         kindVW.writeUint,
		reflect.Uint16:        kindVW.writeUint,
		reflect.Uint32:        kindVW.writeUint,
		reflect.Uint64:        kindVW.writeUint,
		reflect.Uintptr:       kindVW.writeUintptr,
		reflect.Float32:       kindVW.writeFloat,
		reflect.Float64:       kindVW.writeFloat,
		reflect.Complex64:     kindVW.writeComplex,
		reflect.Complex128:    kindVW.writeComplex,
		reflect.Array:         kindVW.writeArray,
		reflect.Chan:          kindVW.writeChan,
		reflect.Func:          kindVW.writeFunc,
		reflect.Interface:     kindVW.writeInterface,
		reflect.Map:           kindVW.writeMap,
		reflect.Pointer:       kindVW.writePointer,
		reflect.Slice:         kindVW.writeSlice,
		reflect.String:        kindVW.writeString,
		reflect.Struct:        kindVW.writeStruct,
		reflect.UnsafePointer: kindVW.writeUnsafePointer,
	}
	return kindVW
}

// WriteValue implements [ValueWriter].
func (vw *KindValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	return vw.ValueWriters[v.Kind()](w, st, v)
}

func (vw *KindValueWriter) writeInvalid(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseInvalid.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeBool(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseBool.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeInt(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseInt.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeUint(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseUint.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeUintptr(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseUintptr.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeFloat(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseFloat.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeComplex(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseComplex.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeArray(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseArray.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeChan(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseChan.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeFunc(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseFunc.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeInterface(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseInterface.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeMap(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseMap.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writePointer(w io.Writer, st State, v reflect.Value) bool {
	return vw.BasePointer.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeSlice(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseSlice.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeString(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseString.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeStruct(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseStruct.WriteValue(w, st, v)
}

func (vw *KindValueWriter) writeUnsafePointer(w io.Writer, st State, v reflect.Value) bool {
	return vw.BaseUnsafePointer.WriteValue(w, st, v)
}
