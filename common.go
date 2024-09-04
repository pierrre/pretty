package pretty

import (
	"reflect"
)

// DefaultCommonValueWriter is the default [CommonValueWriter].
var DefaultCommonValueWriter = NewCommonValueWriter()

// CommonValueWriter is a [ValueWriter] with common [ValueWriter].
//
// Any [ValueWriter] can be configured, but it's not allowed to change the pointer value.
// Any [ValueWriter] can be set to nil in order to disable it.
// It is not allowed to updated the wrapped [ValueWriter].
//
// It should be created with [NewCommonValueWriter].
type CommonValueWriter struct {
	PanicRecover     *PanicRecoverValueWriter
	UnwrapInterface  *UnwrapInterfaceValueWriter
	Recursion        *RecursionValueWriter
	MaxDepth         *MaxDepthValueWriter
	CanInterface     *CanInterfaceValueWriter
	TypeAndValue     *TypeAndValueWriter
	Type             *TypeValueWriter
	ValueWriters     ValueWriters
	ReflectValue     *ReflectValueWriter
	Time             *TimeValueWriter
	Error            *ErrorValueWriter
	BytesHexDump     *BytesHexDumpValueWriter
	BytesableHexDump *BytesableHexDumpValueWriter
	Stringer         *StringerValueWriter
	Kind             *KindValueWriter
}

// NewCommonValueWriter creates a new [CommonValueWriter] initialized with default values.
func NewCommonValueWriter() *CommonValueWriter {
	vw := &CommonValueWriter{}
	vw.PanicRecover = NewPanicRecoverValueWriter(vw.postPanicRecover)
	vw.UnwrapInterface = NewUnwrapInterfaceValueWriter(vw.postUnwrapInterface)
	vw.Recursion = NewRecursionValueWriter(vw.postRecursion)
	vw.MaxDepth = NewMaxDepthValueWriter(vw.postMaxDepth)
	vw.CanInterface = NewCanInterfaceValueWriter(vw.postCanInterface)
	vw.TypeAndValue = NewTypeAndValueWriter(vw.writeType, vw.postTypeAndValue)
	vw.Type = NewTypeValueWriter()
	vw.ReflectValue = NewReflectValueWriter(vw.loopback)
	vw.Time = NewTimeValueWriter()
	vw.Error = NewErrorValueWriter()
	vw.BytesHexDump = NewBytesHexDumpValueWriter()
	vw.BytesableHexDump = NewBytesableHexDumpValueWriter()
	vw.Stringer = NewStringerValueWriter()
	vw.Kind = NewKindValueWriter(vw.loopback)
	return vw
}

// SetShowLen sets ShowLen on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowLen(show bool) {
	vw.Kind.BaseChan.ShowLen = show
	vw.Kind.BaseMap.ShowLen = show
	vw.Kind.BaseSlice.ShowLen = show
	vw.Kind.BaseString.ShowLen = show
	vw.BytesHexDump.ShowLen = show
	vw.BytesableHexDump.ShowLen = show
	vw.Stringer.ShowLen = show
}

// SetShowCap sets ShowCap on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowCap(show bool) {
	vw.Kind.BaseChan.ShowCap = show
	vw.Kind.BaseSlice.ShowCap = show
	vw.BytesHexDump.ShowCap = show
	vw.BytesableHexDump.ShowCap = show
}

// SetShowAddr sets ShowAddr on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowAddr(show bool) {
	vw.Kind.BaseChan.ShowAddr = show
	vw.Kind.BaseFunc.ShowAddr = show
	vw.Kind.BaseMap.ShowAddr = show
	vw.Kind.BasePointer.ShowAddr = show
	vw.Kind.BaseSlice.ShowAddr = show
	vw.Kind.BaseString.ShowAddr = show
	vw.BytesHexDump.ShowAddr = show
	vw.BytesableHexDump.ShowAddr = show
}

// ConfigureTest configures the [CommonValueWriter] for testing.
//
// It makes the result deterministic.
// It sorts the keys of maps and disables the address.
func (vw *CommonValueWriter) ConfigureTest() {
	vw.Kind.BaseMap.SortKeys = true
	vw.SetShowAddr(false)
}

// WriteValue implements [ValueWriter].
func (vw *CommonValueWriter) WriteValue(st *State, v reflect.Value) bool {
	return vw.panicRecover(st, v)
}

func (vw *CommonValueWriter) panicRecover(st *State, v reflect.Value) bool {
	if vw.PanicRecover == nil {
		return vw.postPanicRecover(st, v)
	}
	return vw.PanicRecover.WriteValue(st, v)
}

func (vw *CommonValueWriter) postPanicRecover(st *State, v reflect.Value) bool {
	return vw.unwrapInterface(st, v)
}

func (vw *CommonValueWriter) loopback(st *State, v reflect.Value) bool {
	return vw.unwrapInterface(st, v)
}

func (vw *CommonValueWriter) unwrapInterface(st *State, v reflect.Value) bool {
	if vw.UnwrapInterface == nil {
		return vw.postUnwrapInterface(st, v)
	}
	return vw.UnwrapInterface.WriteValue(st, v)
}

func (vw *CommonValueWriter) postUnwrapInterface(st *State, v reflect.Value) bool {
	return vw.recursion(st, v)
}

func (vw *CommonValueWriter) recursion(st *State, v reflect.Value) bool {
	if vw.Recursion == nil {
		return vw.postRecursion(st, v)
	}
	return vw.Recursion.WriteValue(st, v)
}

func (vw *CommonValueWriter) postRecursion(st *State, v reflect.Value) bool {
	return vw.maxDepth(st, v)
}

func (vw *CommonValueWriter) maxDepth(st *State, v reflect.Value) bool {
	if vw.MaxDepth == nil {
		return vw.postMaxDepth(st, v)
	}
	return vw.MaxDepth.WriteValue(st, v)
}

func (vw *CommonValueWriter) postMaxDepth(st *State, v reflect.Value) bool {
	return vw.canInterface(st, v)
}

func (vw *CommonValueWriter) canInterface(st *State, v reflect.Value) bool {
	if vw.CanInterface == nil {
		return vw.postCanInterface(st, v)
	}
	return vw.CanInterface.WriteValue(st, v)
}

func (vw *CommonValueWriter) postCanInterface(st *State, v reflect.Value) bool {
	return vw.typeAndValue(st, v)
}

func (vw *CommonValueWriter) writeType(st *State, v reflect.Value) bool {
	return vw.Type.WriteValue(st, v)
}

func (vw *CommonValueWriter) typeAndValue(st *State, v reflect.Value) bool {
	if vw.TypeAndValue == nil || vw.Type == nil {
		return vw.postTypeAndValue(st, v)
	}
	return vw.TypeAndValue.WriteValue(st, v)
}

func (vw *CommonValueWriter) postTypeAndValue(st *State, v reflect.Value) bool {
	return vw.internal(st, v)
}

//nolint:gocyclo // We need to call all [ValueWriter].
func (vw *CommonValueWriter) internal(st *State, v reflect.Value) bool {
	if vw.ValueWriters.WriteValue(st, v) {
		return true
	}
	if vw.ReflectValue != nil && vw.ReflectValue.WriteValue(st, v) {
		return true
	}
	if vw.Time != nil && vw.Time.WriteValue(st, v) {
		return true
	}
	if vw.Error != nil && vw.Error.WriteValue(st, v) {
		return true
	}
	if vw.BytesHexDump != nil && vw.BytesHexDump.WriteValue(st, v) {
		return true
	}
	if vw.BytesableHexDump != nil && vw.BytesableHexDump.WriteValue(st, v) {
		return true
	}
	if vw.Stringer != nil && vw.Stringer.WriteValue(st, v) {
		return true
	}
	if vw.Kind != nil && vw.Kind.WriteValue(st, v) {
		return true
	}
	return false
}
