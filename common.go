package pretty

import (
	"reflect"
	"testing"

	"github.com/pierrre/go-libs/reflectutil"
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
// Any [ValueWriter] can be configured, but it's not allowed to change the pointer value.
// Any [ValueWriter] can be set to nil in order to disable it.
// It is not allowed to updated the wrapped [ValueWriter].
//
// It should be created with [NewCommonWriter].
type CommonWriter struct {
	// UnwrapInterface indicates whether to unwrap interface values.
	// Default: true.
	UnwrapInterface bool
	// RecursionCheck indicates whether to check for infinite recursion.
	// Default: true.
	RecursionCheck bool
	// MaxDepth indicates the maximum depth to write.
	// Default: 0 (no limit).
	MaxDepth int
	// CanInterface indicates whether to convert the [reflect.Value] so it can be used with [reflect.Value.Interface].
	// Default: true.
	CanInterface bool
	// ShowType indicates whether to show the type of values.
	// Default: true.
	ShowType bool
	// Type is the [ValueWriter] for types.
	Type TypeWriter

	// The [ValueWriter]s below can be set to nil to disable them.
	ByType           ByTypeWriters
	ValueWriters     ValueWriters
	Support          *SupportWriter
	Time             *TimeWriter
	BytesHexDump     *BytesHexDumpWriter
	MathBigInt       *MathBigIntWriter
	IterSeq          *IterSeqWriter
	IterSeq2         *IterSeq2Writer
	WeakPointer      *WeakPointerWriter
	ReflectValue     *ReflectValueWriter
	ReflectType      *ReflectTypeWriter
	Error            *ErrorWriter
	BytesableHexDump *BytesableHexDumpWriter
	Stringer         *StringerWriter

	// Kind is the default [ValueWriter].
	// It must not be set to nil.
	Kind *KindWriter
}

// NewCommonWriter creates a new [CommonWriter] initialized with default values.
func NewCommonWriter() *CommonWriter {
	vw := &CommonWriter{}
	vw.UnwrapInterface = true
	vw.RecursionCheck = true
	vw.MaxDepth = 0
	vw.CanInterface = true
	vw.ShowType = true
	vw.Type = *NewTypeWriter(ValueWriterFunc(vw.writeValue))
	vw.ByType = NewByTypeWriters()
	vw.Support = NewSupportWriter()
	vw.Support.Checkers = []SupportChecker{
		SupportCheckerFunc(vw.supports),
	}
	vw.Time = NewTimeWriter()
	vw.BytesHexDump = NewBytesHexDumpWriter()
	vw.MathBigInt = NewMathBigIntWriter()
	vw.IterSeq = NewIterSeqWriter(vw)
	vw.IterSeq2 = NewIterSeq2Writer(vw)
	vw.WeakPointer = NewWeakPointerWriter(vw)
	vw.ReflectValue = NewReflectValueWriter(vw)
	vw.ReflectType = NewReflectTypeWriter()
	vw.Error = NewErrorWriter()
	vw.BytesableHexDump = NewBytesableHexDumpWriter()
	vw.Stringer = NewStringerWriter()
	vw.Kind = NewKindWriter(vw)
	return vw
}

// SetShowLen sets ShowLen on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowLen(show bool) {
	vw.Kind.Chan.ShowLen = show
	vw.Kind.Map.ShowLen = show
	vw.Kind.Slice.ShowLen = show
	vw.Kind.String.ShowLen = show
	vw.BytesHexDump.ShowLen = show
	vw.BytesableHexDump.ShowLen = show
	vw.Stringer.ShowLen = show
}

// SetShowCap sets ShowCap on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowCap(show bool) {
	vw.Kind.Chan.ShowCap = show
	vw.Kind.Slice.ShowCap = show
	vw.BytesHexDump.ShowCap = show
	vw.BytesableHexDump.ShowCap = show
}

// SetShowAddr sets ShowAddr on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowAddr(show bool) {
	vw.Kind.Chan.ShowAddr = show
	vw.Kind.Func.ShowAddr = show
	vw.Kind.Map.ShowAddr = show
	vw.Kind.Pointer.ShowAddr = show
	vw.Kind.Slice.ShowAddr = show
	vw.Kind.String.ShowAddr = show
	vw.BytesHexDump.ShowAddr = show
	vw.BytesableHexDump.ShowAddr = show
}

// SetShowIndexes sets ShowIndexes on all [ValueWriter] that supports it.
func (vw *CommonWriter) SetShowIndexes(show bool) {
	vw.Kind.Array.ShowIndexes = show
	vw.Kind.Slice.ShowIndexes = show
	vw.Kind.Chan.ShowIndexes = show
}

// ConfigureTest configures the [CommonWriter] for testing.
//
// It makes the result deterministic.
// It sorts the keys of maps and disables the address/capacity.
// The enabled boolean is used to enable or disable the configuration.
func (vw *CommonWriter) ConfigureTest(enabled bool) {
	vw.Kind.Map.SortKeys = enabled
	vw.SetShowAddr(!enabled)
	vw.SetShowCap(!enabled)
}

// WriteValue implements [ValueWriter].
//
//nolint:gocyclo // Yes it's complex.
func (vw *CommonWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.UnwrapInterface {
		var isNil bool
		v, isNil = unwrapInterface(st, v)
		if isNil {
			return true
		}
	}
	if vw.RecursionCheck {
		showInfos := vw.Kind != nil && vw.Kind.Pointer.ShowAddr
		visitedAdded, recursionDetected := checkRecursion(st, v, showInfos)
		if recursionDetected {
			return true
		}
		if visitedAdded {
			defer postRecursion(st)
		}
	}
	if vw.MaxDepth > 0 {
		if checkMaxDepth(st, vw.MaxDepth) {
			return true
		}
		defer postMaxDepth(st)
	}
	if vw.CanInterface {
		v, _ = reflectutil.ConvertValueCanInterface(v)
	}
	if vw.ShowType {
		return vw.Type.WriteValue(st, v)
	}
	return vw.writeValue(st, v)
}

//nolint:gocyclo // We need to call all [ValueWriter].
func (vw *CommonWriter) writeValue(st *State, v reflect.Value) bool {
	if vw.ByType != nil && vw.ByType.WriteValue(st, v) {
		return true
	}
	if vw.ValueWriters != nil && vw.ValueWriters.WriteValue(st, v) {
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
	if vw.MathBigInt != nil && vw.MathBigInt.WriteValue(st, v) {
		return true
	}
	if vw.IterSeq != nil && vw.IterSeq.WriteValue(st, v) {
		return true
	}
	if vw.IterSeq2 != nil && vw.IterSeq2.WriteValue(st, v) {
		return true
	}
	if vw.WeakPointer != nil && vw.WeakPointer.WriteValue(st, v) {
		return true
	}
	if vw.ReflectValue != nil && vw.ReflectValue.WriteValue(st, v) {
		return true
	}
	if vw.ReflectType != nil && vw.ReflectType.WriteValue(st, v) {
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

//nolint:gocyclo // We need to call all [SupportChecker].
func (vw *CommonWriter) supports(typ reflect.Type) ValueWriter {
	if f := callSupportCheckerPointer(vw.Time, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.BytesHexDump, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.MathBigInt, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.IterSeq, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.IterSeq2, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.WeakPointer, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.ReflectValue, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.ReflectType, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.Error, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.BytesableHexDump, typ); f != nil {
		return f
	}
	if f := callSupportCheckerPointer(vw.Stringer, typ); f != nil {
		return f
	}
	return vw.Kind.Supports(typ)
}

func callSupportCheckerPointer[P interface {
	*T
	SupportChecker
}, T any](p P, typ reflect.Type) ValueWriter {
	if p != nil {
		return p.Supports(typ)
	}
	return nil
}
