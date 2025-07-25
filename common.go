package pretty

import (
	"reflect"
	"testing"
)

// DefaultCommonWriter is the default [CommonWriter].
//
// It is configured with [CommonWriter.ConfigureTest] and [testing.Testing].
var DefaultCommonWriter = NewCommonWriter()

func init() {
	DefaultCommonWriter.ConfigureTest(testing.Testing())
}

// CommonWriter is a [ValueWriter] with common [ValueWriter].
//
// Any [ValueWriter] can be configured, and can be set to nil in order to disable it.
//
// It should be created with [NewCommonWriter].
type CommonWriter struct {
	UnwrapInterface  *UnwrapInterfaceWriter
	Recursion        *RecursionWriter
	MaxDepth         *MaxDepthWriter
	CanInterface     *CanInterfaceWriter
	Type             *TypeWriter
	ByType           ByTypeWriters
	ValueWriters     ValueWriters
	Support          *SupportWriter
	Time             *TimeWriter
	BytesHexDump     *BytesHexDumpWriter
	MathBig          *MathBigWriter
	Iter             *IterWriter
	Range            *RangeWriter
	WeakPointer      *WeakPointerWriter
	Reflect          *ReflectWriter
	Error            *ErrorWriter
	BytesableHexDump *BytesableHexDumpWriter
	Stringer         *StringerWriter
	Kind             *KindWriter
}

// NewCommonWriter creates a new [CommonWriter] initialized with default values.
func NewCommonWriter() *CommonWriter {
	vw := &CommonWriter{}
	vw.UnwrapInterface = NewUnwrapInterfaceWriter(nil)
	vw.Recursion = NewRecursionWriter(nil)
	vw.MaxDepth = NewMaxDepthWriter(nil)
	vw.CanInterface = NewCanInterfaceWriter(nil)
	vw.Type = NewTypeWriter(nil)
	vw.ByType = NewByTypeWriters()
	vw.Support = NewSupportWriter()
	vw.Support.Checkers = []SupportChecker{
		vw,
	}
	vw.Time = NewTimeWriter()
	vw.BytesHexDump = NewBytesHexDumpWriter()
	vw.MathBig = NewMathBigWriter()
	vw.Iter = NewIterWriter(vw)
	vw.Range = NewRangeWriter(vw)
	vw.WeakPointer = NewWeakPointerWriter(vw)
	vw.Reflect = NewReflectWriter(vw)
	vw.Error = NewErrorWriter(vw)
	vw.BytesableHexDump = NewBytesableHexDumpWriter()
	vw.Stringer = NewStringerWriter()
	vw.Kind = NewKindWriter(vw)
	return vw
}

// SetShowLen sets ShowLen on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowLen(show bool) {
	if vw.Kind != nil {
		vw.Kind.Chan.ShowLen = show
		vw.Kind.Map.ShowLen = show
		vw.Kind.Slice.ShowLen = show
		vw.Kind.String.ShowLen = show
	}
	if vw.BytesHexDump != nil {
		vw.BytesHexDump.ShowLen = show
	}
	if vw.BytesableHexDump != nil {
		vw.BytesableHexDump.ShowLen = show
	}
	if vw.Stringer != nil {
		vw.Stringer.ShowLen = show
	}
}

// SetShowCap sets ShowCap on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowCap(show bool) {
	if vw.Kind != nil {
		vw.Kind.Chan.ShowCap = show
		vw.Kind.Slice.ShowCap = show
	}
	if vw.BytesHexDump != nil {
		vw.BytesHexDump.ShowCap = show
	}
	if vw.BytesableHexDump != nil {
		vw.BytesableHexDump.ShowCap = show
	}
}

// SetShowAddr sets ShowAddr on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowAddr(show bool) {
	if vw.Recursion != nil {
		vw.Recursion.ShowAddr = show
	}
	if vw.Kind != nil {
		vw.Kind.Chan.ShowAddr = show
		vw.Kind.Func.ShowAddr = show
		vw.Kind.Map.ShowAddr = show
		vw.Kind.Pointer.ShowAddr = show
		vw.Kind.Slice.ShowAddr = show
		vw.Kind.String.ShowAddr = show
	}
	if vw.BytesHexDump != nil {
		vw.BytesHexDump.ShowAddr = show
	}
	if vw.BytesableHexDump != nil {
		vw.BytesableHexDump.ShowAddr = show
	}
}

// SetShowIndexes sets ShowIndexes on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowIndexes(show bool) {
	if vw.Kind != nil {
		vw.Kind.Array.ShowIndexes = show
		vw.Kind.Slice.ShowIndexes = show
		vw.Kind.Chan.ShowIndexes = show
	}
}

// ConfigureTest configures the [CommonWriter] for testing.
//
// It makes the result deterministic.
// It sorts the keys of maps and disables the address/capacity.
// The enabled boolean is used to enable or disable the configuration.
func (vw *CommonWriter) ConfigureTest(enabled bool) {
	if vw.Kind != nil {
		vw.Kind.Map.SortKeys = enabled
	}
	vw.SetShowAddr(!enabled)
	vw.SetShowCap(!enabled)
}

// WriteValue implements [ValueWriter].
func (vw *CommonWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.UnwrapInterface != nil {
		var isNil bool
		v, isNil = vw.UnwrapInterface.unwrapInterface(st, v)
		if isNil {
			return true
		}
	}
	if vw.Recursion != nil {
		recursionVisitedAdded, recursionDetected := vw.Recursion.checkRecursion(st, v)
		if recursionDetected {
			return true
		}
		if recursionVisitedAdded {
			defer vw.Recursion.postRecursion(st)
		}
	}
	if vw.MaxDepth != nil {
		maxReached := vw.MaxDepth.checkMaxDepth(st)
		defer vw.MaxDepth.postMaxDepth(st)
		if maxReached {
			return true
		}
	}
	if vw.CanInterface != nil {
		v = vw.CanInterface.convertValue(v)
	}
	if vw.Type != nil {
		knownType := vw.Type.writeType(st, v)
		defer vw.Type.postType(st, knownType)
	}
	return vw.writeValue(st, v)
}

//nolint:gocyclo // We need to call all [ValueWriter].
func (vw *CommonWriter) writeValue(st *State, v reflect.Value) bool {
	if len(vw.ByType) != 0 && vw.ByType.WriteValue(st, v) {
		return true
	}
	if len(vw.ValueWriters) != 0 && vw.ValueWriters.WriteValue(st, v) {
		return true
	}
	if vw.Support != nil && vw.Support.WriteValue(st, v) {
		return true
	}
	if vw.Time != nil && vw.Time.WriteValue(st, v) {
		return true
	}
	if vw.BytesHexDump != nil && vw.BytesHexDump.WriteValue(st, v) {
		return true
	}
	if vw.MathBig != nil && vw.MathBig.WriteValue(st, v) {
		return true
	}
	if vw.Iter != nil && vw.Iter.WriteValue(st, v) {
		return true
	}
	if vw.Range != nil && vw.Range.WriteValue(st, v) {
		return true
	}
	if vw.WeakPointer != nil && vw.WeakPointer.WriteValue(st, v) {
		return true
	}
	if vw.Reflect != nil && vw.Reflect.WriteValue(st, v) {
		return true
	}
	if vw.Error != nil && vw.Error.WriteValue(st, v) {
		return true
	}
	if vw.BytesableHexDump != nil && vw.BytesableHexDump.WriteValue(st, v) {
		return true
	}
	if vw.Stringer != nil && vw.Stringer.WriteValue(st, v) {
		return true
	}
	return vw.Kind.WriteValue(st, v)
}

// Supports implements [SupportChecker].
//
//nolint:gocyclo // We need to call all [SupportChecker].
func (vw *CommonWriter) Supports(typ reflect.Type) ValueWriter {
	if w := callSupportCheckerPointer(vw.Time, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.BytesHexDump, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.MathBig, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Iter, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Range, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.WeakPointer, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Reflect, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Error, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.BytesableHexDump, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Stringer, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Kind, typ); w != nil {
		return w
	}
	return nil
}
