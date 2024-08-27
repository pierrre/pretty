package pretty

import (
	"io"
	"reflect"
)

// KindValueWriter is a [ValueWriter] that writes the value with the kind-specific [ValueWriter].
//
// The "Base*" fields are the default [ValueWriter] for each group of kinds.
// They can be configured.
//
// The other fields are the [ValueWriter] for each kind.
// It's allowed to update theses fields.
// Each [ValueWriter] must be able to handle the corresponding kind.
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

	Invalid       ValueWriter
	Bool          ValueWriter
	Int           ValueWriter
	Int8          ValueWriter
	Int16         ValueWriter
	Int32         ValueWriter
	Int64         ValueWriter
	Uint          ValueWriter
	Uint8         ValueWriter
	Uint16        ValueWriter
	Uint32        ValueWriter
	Uint64        ValueWriter
	Uintptr       ValueWriter
	Float32       ValueWriter
	Float64       ValueWriter
	Complex64     ValueWriter
	Complex128    ValueWriter
	Array         ValueWriter
	Chan          ValueWriter
	Func          ValueWriter
	Interface     ValueWriter
	Map           ValueWriter
	Pointer       ValueWriter
	Slice         ValueWriter
	String        ValueWriter
	Struct        ValueWriter
	UnsafePointer ValueWriter
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
	kindVW.Invalid = kindVW.writeInvalid
	kindVW.Bool = kindVW.writeBool
	kindVW.Int = kindVW.writeInt
	kindVW.Int8 = kindVW.writeInt
	kindVW.Int16 = kindVW.writeInt
	kindVW.Int32 = kindVW.writeInt
	kindVW.Int64 = kindVW.writeInt
	kindVW.Uint = kindVW.writeUint
	kindVW.Uint8 = kindVW.writeUint
	kindVW.Uint16 = kindVW.writeUint
	kindVW.Uint32 = kindVW.writeUint
	kindVW.Uint64 = kindVW.writeUint
	kindVW.Uintptr = kindVW.writeUintptr
	kindVW.Float32 = kindVW.writeFloat
	kindVW.Float64 = kindVW.writeFloat
	kindVW.Complex64 = kindVW.writeComplex
	kindVW.Complex128 = kindVW.writeComplex
	kindVW.Array = kindVW.writeArray
	kindVW.Chan = kindVW.writeChan
	kindVW.Func = kindVW.writeFunc
	kindVW.Interface = kindVW.writeInterface
	kindVW.Map = kindVW.writeMap
	kindVW.Pointer = kindVW.writePointer
	kindVW.Slice = kindVW.writeSlice
	kindVW.String = kindVW.writeString
	kindVW.Struct = kindVW.writeStruct
	kindVW.UnsafePointer = kindVW.writeUnsafePointer
	return kindVW
}

// WriteValue implements [ValueWriter].
//
//nolint:gocyclo // We need to handle all kinds.
func (vw *KindValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // All kinds are handled, "invalid" is in the default case.
	case reflect.Bool:
		mustHandle(vw.Bool(w, st, v))
	case reflect.Int:
		mustHandle(vw.Int(w, st, v))
	case reflect.Int8:
		mustHandle(vw.Int8(w, st, v))
	case reflect.Int16:
		mustHandle(vw.Int16(w, st, v))
	case reflect.Int32:
		mustHandle(vw.Int32(w, st, v))
	case reflect.Int64:
		mustHandle(vw.Int64(w, st, v))
	case reflect.Uint:
		mustHandle(vw.Uint(w, st, v))
	case reflect.Uint8:
		mustHandle(vw.Uint8(w, st, v))
	case reflect.Uint16:
		mustHandle(vw.Uint16(w, st, v))
	case reflect.Uint32:
		mustHandle(vw.Uint32(w, st, v))
	case reflect.Uint64:
		mustHandle(vw.Uint64(w, st, v))
	case reflect.Uintptr:
		mustHandle(vw.Uintptr(w, st, v))
	case reflect.Float32:
		mustHandle(vw.Float32(w, st, v))
	case reflect.Float64:
		mustHandle(vw.Float64(w, st, v))
	case reflect.Complex64:
		mustHandle(vw.Complex64(w, st, v))
	case reflect.Complex128:
		mustHandle(vw.Complex128(w, st, v))
	case reflect.Array:
		mustHandle(vw.Array(w, st, v))
	case reflect.Chan:
		mustHandle(vw.Chan(w, st, v))
	case reflect.Func:
		mustHandle(vw.Func(w, st, v))
	case reflect.Interface:
		mustHandle(vw.Interface(w, st, v))
	case reflect.Map:
		mustHandle(vw.Map(w, st, v))
	case reflect.Pointer:
		mustHandle(vw.Pointer(w, st, v))
	case reflect.Slice:
		mustHandle(vw.Slice(w, st, v))
	case reflect.String:
		mustHandle(vw.String(w, st, v))
	case reflect.Struct:
		mustHandle(vw.Struct(w, st, v))
	case reflect.UnsafePointer:
		mustHandle(vw.UnsafePointer(w, st, v))
	default:
		mustHandle(vw.Invalid(w, st, v))
	}
	return true
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
