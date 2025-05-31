package pretty

import (
	"reflect"
	"testing"

	"github.com/pierrre/go-libs/reflectutil"
)

// DefaultCommonValueWriter is the default [CommonValueWriter].
//
// It is configured with [CommonValueWriter.ConfigureTest] and [testing.Testing].
var DefaultCommonValueWriter = NewCommonValueWriter()

func init() {
	DefaultCommonValueWriter.ConfigureTest(testing.Testing())
}

// CommonValueWriter is a [ValueWriter] with common [ValueWriter].
//
// Any [ValueWriter] can be configured, but it's not allowed to change the pointer value.
// Any [ValueWriter] can be set to nil in order to disable it.
// It is not allowed to updated the wrapped [ValueWriter].
//
// It should be created with [NewCommonValueWriter].
type CommonValueWriter struct {
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
	Type TypeValueWriter

	// THe [ValueWriter]s below can be set to nil to disable them.
	ByTypeValueWriters ByTypeValueWriters
	ValueWriters       ValueWriters
	ReflectValue       *ReflectValueWriter
	ReflectType        *ReflectTypeWriter
	Time               *TimeValueWriter
	Error              *ErrorValueWriter
	MathBigInt         *MathBigIntWriter
	WeakPointer        *WeakPointerWriter
	BytesHexDump       *BytesHexDumpValueWriter
	BytesableHexDump   *BytesableHexDumpValueWriter
	Stringer           *StringerValueWriter
	Iter               *IterValueWriter
	Kind               *KindValueWriter
}

// NewCommonValueWriter creates a new [CommonValueWriter] initialized with default values.
func NewCommonValueWriter() *CommonValueWriter {
	vw := &CommonValueWriter{}
	vw.UnwrapInterface = true
	vw.RecursionCheck = true
	vw.MaxDepth = 0
	vw.CanInterface = true
	vw.ShowType = true
	vw.Type = *NewTypeValueWriter(ValueWriterFunc(vw.internal))
	vw.ByTypeValueWriters = NewByTypeValueWriters()
	vw.ReflectValue = NewReflectValueWriter(vw)
	vw.ReflectType = NewReflectTypeWriter()
	vw.Time = NewTimeValueWriter()
	vw.Error = NewErrorValueWriter()
	vw.MathBigInt = NewMathBigIntWriter()
	vw.WeakPointer = NewWeakPointerWriter(vw)
	vw.BytesHexDump = NewBytesHexDumpValueWriter()
	vw.BytesableHexDump = NewBytesableHexDumpValueWriter()
	vw.Stringer = NewStringerValueWriter()
	vw.Iter = NewIterValueWriter(vw)
	vw.Kind = NewKindValueWriter(vw)
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

// SetShowIndexes sets ShowIndexes on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowIndexes(show bool) {
	vw.Kind.BaseArray.ShowIndexes = show
	vw.Kind.BaseSlice.ShowIndexes = show
	vw.Kind.BaseChan.ShowIndexes = show
}

// ConfigureTest configures the [CommonValueWriter] for testing.
//
// It makes the result deterministic.
// It sorts the keys of maps and disables the address/capacity.
// The enabled boolean is used to enable or disable the configuration.
func (vw *CommonValueWriter) ConfigureTest(enabled bool) {
	vw.Kind.BaseMap.SortKeys = enabled
	vw.SetShowAddr(!enabled)
	vw.SetShowCap(!enabled)
}

// WriteValue implements [ValueWriter].
func (vw *CommonValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.UnwrapInterface {
		var isNil bool
		v, isNil = unwrapInterface(st, v)
		if isNil {
			return true
		}
	}
	if vw.RecursionCheck {
		visitedAdded, recursionDetected := checkRecursion(st, v)
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
	return vw.internal(st, v)
}

//nolint:gocyclo // We need to call all [ValueWriter].
func (vw *CommonValueWriter) internal(st *State, v reflect.Value) bool {
	if vw.ByTypeValueWriters != nil && vw.ByTypeValueWriters.WriteValue(st, v) {
		return true
	}
	if vw.ValueWriters.WriteValue(st, v) {
		return true
	}
	if vw.ReflectValue != nil && vw.ReflectValue.WriteValue(st, v) {
		return true
	}
	if vw.ReflectType != nil && vw.ReflectType.WriteValue(st, v) {
		return true
	}
	if vw.Time != nil && vw.Time.WriteValue(st, v) {
		return true
	}
	if vw.Error != nil && vw.Error.WriteValue(st, v) {
		return true
	}
	if vw.MathBigInt != nil && vw.MathBigInt.WriteValue(st, v) {
		return true
	}
	if vw.WeakPointer != nil && vw.WeakPointer.WriteValue(st, v) {
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
	if vw.Iter != nil && vw.Iter.WriteValue(st, v) {
		return true
	}
	if vw.Kind != nil && vw.Kind.WriteValue(st, v) {
		return true
	}
	return false
}
