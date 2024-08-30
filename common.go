package pretty

import (
	"io"
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

// ConfigureWithPrinter configures the [CommonValueWriter] with a [Printer].
//
// It sets the [MapValueWriter.SortKeysCmpDefault] to use [Printer.compare].
func (vw *CommonValueWriter) ConfigureWithPrinter(p *Printer) {
	vw.Kind.BaseMap.SortKeysCmpDefault = p.compare
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
func (vw *CommonValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	return vw.panicRecover(w, st, v)
}

func (vw *CommonValueWriter) panicRecover(w io.Writer, st State, v reflect.Value) bool {
	if vw.PanicRecover == nil {
		return vw.postPanicRecover(w, st, v)
	}
	return vw.PanicRecover.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) postPanicRecover(w io.Writer, st State, v reflect.Value) bool {
	return vw.unwrapInterface(w, st, v)
}

func (vw *CommonValueWriter) loopback(w io.Writer, st State, v reflect.Value) bool {
	return vw.unwrapInterface(w, st, v)
}

func (vw *CommonValueWriter) unwrapInterface(w io.Writer, st State, v reflect.Value) bool {
	if vw.UnwrapInterface == nil {
		return vw.postUnwrapInterface(w, st, v)
	}
	return vw.UnwrapInterface.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) postUnwrapInterface(w io.Writer, st State, v reflect.Value) bool {
	return vw.recursion(w, st, v)
}

func (vw *CommonValueWriter) recursion(w io.Writer, st State, v reflect.Value) bool {
	if vw.Recursion == nil {
		return vw.postRecursion(w, st, v)
	}
	return vw.Recursion.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) postRecursion(w io.Writer, st State, v reflect.Value) bool {
	return vw.maxDepth(w, st, v)
}

func (vw *CommonValueWriter) maxDepth(w io.Writer, st State, v reflect.Value) bool {
	if vw.MaxDepth == nil {
		return vw.postMaxDepth(w, st, v)
	}
	return vw.MaxDepth.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) postMaxDepth(w io.Writer, st State, v reflect.Value) bool {
	return vw.canInterface(w, st, v)
}

func (vw *CommonValueWriter) canInterface(w io.Writer, st State, v reflect.Value) bool {
	if vw.CanInterface == nil {
		return vw.postCanInterface(w, st, v)
	}
	return vw.CanInterface.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) postCanInterface(w io.Writer, st State, v reflect.Value) bool {
	return vw.typeAndValue(w, st, v)
}

func (vw *CommonValueWriter) writeType(w io.Writer, st State, v reflect.Value) bool {
	return vw.Type.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) typeAndValue(w io.Writer, st State, v reflect.Value) bool {
	if vw.TypeAndValue == nil || vw.Type == nil {
		return vw.postTypeAndValue(w, st, v)
	}
	return vw.TypeAndValue.WriteValue(w, st, v)
}

func (vw *CommonValueWriter) postTypeAndValue(w io.Writer, st State, v reflect.Value) bool {
	return vw.internal(w, st, v)
}

//nolint:gocyclo // We need to call all [ValueWriter].
func (vw *CommonValueWriter) internal(w io.Writer, st State, v reflect.Value) bool {
	if vw.ValueWriters.WriteValue(w, st, v) {
		return true
	}
	if vw.ReflectValue != nil && vw.ReflectValue.WriteValue(w, st, v) {
		return true
	}
	if vw.Time != nil && vw.Time.WriteValue(w, st, v) {
		return true
	}
	if vw.Error != nil && vw.Error.WriteValue(w, st, v) {
		return true
	}
	if vw.BytesHexDump != nil && vw.BytesHexDump.WriteValue(w, st, v) {
		return true
	}
	if vw.BytesableHexDump != nil && vw.BytesableHexDump.WriteValue(w, st, v) {
		return true
	}
	if vw.Stringer != nil && vw.Stringer.WriteValue(w, st, v) {
		return true
	}
	if vw.Kind != nil && vw.Kind.WriteValue(w, st, v) {
		return true
	}
	return false
}
