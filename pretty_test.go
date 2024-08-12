package pretty_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe" //nolint:depguard // Required for test.

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	"github.com/pierrre/go-libs/reflectutil"
	. "github.com/pierrre/pretty"
)

func init() {
	DefaultCommonValueWriter.ConfigureTest()
}

func Example() {
	type exampleStruct struct {
		Int    int
		Float  float64
		String string
		Map    map[string]int
		Slice  []int
	}
	v := exampleStruct{
		Int:    123,
		Float:  123.456,
		String: "test",
		Map: map[string]int{
			"foo": 1,
			"bar": 2,
		},
		Slice: []int{1, 2, 3},
	}
	s := String(v)
	fmt.Println(s)
	// Output:
	// [pretty_test.exampleStruct] {
	// 	Int: [int] 123,
	// 	Float: [float64] 123.456,
	// 	String: [string] (len=4) "test",
	// 	Map: [map[string]int] (len=2) {
	// 		(len=3) "bar": 2,
	// 		(len=3) "foo": 1,
	// 	},
	// 	Slice: [[]int] (len=3 cap=3) {
	// 		1,
	// 		2,
	// 		3,
	// 	},
	// }
}

func ExampleString() {
	s := String("test")
	fmt.Println(s)
	// Output: [string] (len=4) "test"
}

func ExampleWrite() {
	buf := new(bytes.Buffer)
	Write(buf, "test")
	s := buf.String()
	fmt.Println(s)
	// Output: [string] (len=4) "test"
}

func ExampleFormatter() {
	f := Formatter("test")
	s := fmt.Sprintf("%v", f)
	fmt.Println(s)
	// Output: [string] (len=4) "test"
}

func ExampleValueWriter() {
	c := NewConfig()
	vw := func(c *Config, w io.Writer, st State, v reflect.Value) bool {
		_, _ = io.WriteString(w, "example")
		return true
	}
	p := NewPrinter(c, vw)
	s := p.String("test")
	fmt.Println(s)
	// Output: example
}

func newTestPrinter() (*Printer, *CommonValueWriter) {
	c := NewConfig()
	vw := NewCommonValueWriter()
	vw.ConfigureTest()
	p := NewPrinterCommon(c, vw)
	return p, vw
}

type testCase struct {
	name         string
	value        any
	panicRecover bool
	configure    func(vw *CommonValueWriter)
	ignoreResult bool
	ignoreAllocs bool
}

func (tc testCase) newPrinter() *Printer {
	p, vw := newTestPrinter()
	if !tc.panicRecover {
		vw.PanicRecover = nil
	}
	if tc.configure != nil {
		tc.configure(vw)
	}
	return p
}

var testCases = []testCase{
	{
		name:  "Nil",
		value: nil,
	},
	{
		name:  "Invalid",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{func(c *Config, w io.Writer, st State, v reflect.Value) bool {
				return vw.Kind.WriteValue(c, w, st, reflect.ValueOf(nil))
			}}
		},
	},
	{
		name:         "KindDisabled",
		value:        "test",
		panicRecover: true,
		configure: func(vw *CommonValueWriter) {
			vw.Kind = nil
		},
	},
	{
		name:  "Bool",
		value: true,
	},
	{
		name:  "BoolNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseBool.WriteValue}
		},
	},
	{
		name:  "Int",
		value: 123,
	},
	{
		name:  "Int8",
		value: int8(123),
	},
	{
		name:  "Int16",
		value: int16(123),
	},
	{
		name:  "Int32",
		value: int32(123),
	},
	{
		name:  "Int64",
		value: int64(123),
	},
	{
		name:  "IntNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseInt.WriteValue}
		},
	},
	{
		name:  "Uint",
		value: uint(123),
	},
	{
		name:  "Uint8",
		value: uint8(123),
	},
	{
		name:  "Uint16",
		value: uint16(123),
	},
	{
		name:  "Uint32",
		value: uint32(123),
	},
	{
		name:  "Uint64",
		value: uint64(123),
	},
	{
		name:  "UintNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseUint.WriteValue}
		},
	},
	{
		name:  "Uintptr",
		value: uintptr(123),
	},
	{
		name:  "UintptrNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseUintptr.WriteValue}
		},
	},
	{
		name:  "Float32",
		value: float32(123.456),
	},
	{
		name:  "Float64",
		value: float64(123.456),
	},
	{
		name:  "FloatNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseFloat.WriteValue}
		},
	},
	{
		name:  "Complex64",
		value: complex64(123.456 + 789.123i),
	},
	{
		name:  "Complex128",
		value: complex128(123.456 + 789.123i),
	},
	{
		name:  "ComplexNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseComplex.WriteValue}
		},
	},
	{
		name:  "Array",
		value: [3]int{1, 2, 3},
	},
	{
		name:  "ArrayEmpty",
		value: [0]int{},
	},
	{
		name:  "ArrayUnknownType",
		value: [3]any{1, 2, 3},
	},
	{
		name:  "ArrayShowKnownTypes",
		value: [3]int{1, 2, 3},
		configure: func(vw *CommonValueWriter) {
			vw.TypeAndValue.ShowKnownTypes = true
		},
	},
	{
		name:  "ArrayNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseArray.WriteValue}
		},
	},
	{
		name: "Chan",
		value: func() chan int {
			c := make(chan int, 5)
			c <- 123
			return c
		}(),
	},
	{
		name:  "ChanNil",
		value: chan int(nil),
	},
	{
		name: "ChanShowAddr",
		value: func() chan int {
			c := make(chan int, 5)
			c <- 123
			return c
		}(),
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseChan.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name:  "ChanNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseChan.WriteValue}
		},
	},
	{
		name:  "Func",
		value: NewConfig,
	},
	{
		name:  "FuncNil",
		value: (func())(nil),
	},
	{
		name:  "FuncShowAddr",
		value: NewConfig,
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseFunc.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name:  "FuncNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseFunc.WriteValue}
		},
	},
	{
		name:  "Interface",
		value: [1]any{123},
		configure: func(vw *CommonValueWriter) {
			vw.UnwrapInterface = nil
		},
	},
	{
		name:  "InterfaceNil",
		value: [1]any{nil},
		configure: func(vw *CommonValueWriter) {
			vw.UnwrapInterface = nil
		},
	},
	{
		name:  "InterfaceNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseInterface.WriteValue}
		},
	},
	{
		name:  "MapNil",
		value: map[int]int(nil),
	},
	{
		name:  "MapEmpty",
		value: map[int]int{},
	},
	{
		name:  "MapShowAddr",
		value: map[int]int{1: 2},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = false
			vw.Kind.BaseMap.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name:  "MapUnsortedExported",
		value: map[int]int{1: 2},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = false
		},
	},
	{
		name:  "MapUnsortedExportedTruncated",
		value: map[int]int{1: 2, 3: 4},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = false
			vw.Kind.BaseMap.MaxLen = 1
		},
		ignoreResult: true,
	},
	{
		name:  "MapUnsortedExportedInterface",
		value: map[any]any{1: 2},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = false
		},
	},
	{
		name:  "MapUnsortedUnexported",
		value: testUnexported{v: map[int]int{1: 2}},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = false
		},
	},
	{
		name:  "MapUnsortedUnexportedTruncated",
		value: testUnexported{v: map[int]int{1: 2, 3: 4}},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = false
			vw.Kind.BaseMap.MaxLen = 1
		},
		ignoreResult: true,
	},
	{
		name:  "MapSortedBool",
		value: map[bool]int{false: 1, true: 2},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
		},
	},
	{
		name:  "MapSortedInt",
		value: map[int]int{1: 2, 3: 4, 5: 6},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
		},
	},
	{
		name:  "MapSortedUint",
		value: map[uint]int{1: 2, 3: 4, 5: 6},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
		},
	},
	{
		name:  "MapSortedFloat",
		value: map[float64]int{1: 2, 3: 4, 5: 6},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
		},
	},
	{
		name:  "MapSortedString",
		value: map[string]int{"a": 1, "b": 2, "c": 3},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
		},
	},
	{
		name:  "MapSortedDefault",
		value: map[testComparableStruct]int{{V: 1}: 2, {V: 3}: 4, {V: 5}: 6},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
		},
	},
	{
		name:  "MapSortedDefaultSimple",
		value: map[testComparableStruct]int{{V: 1}: 2, {V: 3}: 4, {V: 5}: 6},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
			vw.Kind.BaseMap.SortKeysCmpDefault = nil
		},
		ignoreAllocs: true,
	},
	{
		name:  "MapSortedTruncated",
		value: map[int]int{1: 2, 3: 4, 5: 6},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseMap.SortKeys = true
			vw.Kind.BaseMap.MaxLen = 2
		},
	},
	{
		name:  "MapUnknownType",
		value: map[any]any{"a": 1, "b": 2, "c": 3},
	},
	{
		name:  "MapShowKnownTypes",
		value: map[string]int{"a": 1, "b": 2, "c": 3},
		configure: func(vw *CommonValueWriter) {
			vw.TypeAndValue.ShowKnownTypes = true
		},
	},
	{
		name:  "MapNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseMap.WriteValue}
		},
	},
	{
		name: "Pointer",
		value: func() *int {
			i := 123
			return &i
		}(),
	},
	{
		name:  "PointerNil",
		value: (*int)(nil),
	},
	{
		name: "PointerShowAddr",
		value: func() *int {
			i := 123
			return &i
		}(),
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BasePointer.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name: "PointerUnknownType",
		value: func() *any {
			i := any(123)
			return &i
		}(),
	},
	{
		name: "PointerShowKnownTypes",
		value: func() *int {
			i := 123
			return &i
		}(),
		configure: func(vw *CommonValueWriter) {
			vw.TypeAndValue.ShowKnownTypes = true
		},
	},
	{
		name:  "PointerNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BasePointer.WriteValue}
		},
	},
	{
		name:  "Slice",
		value: []int{1, 2, 3},
	},
	{
		name:  "SliceNil",
		value: []int(nil),
	},
	{
		name:  "SliceShowAddr",
		value: []int{1, 2, 3},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseSlice.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name:  "SliceEmpty",
		value: []int{},
	},
	{
		name:  "SliceTruncated",
		value: []int{1, 2, 3},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseSlice.MaxLen = 2
		},
	},
	{
		name:  "SliceUnknownType",
		value: []any{1, 2, 3},
	},
	{
		name:  "SliceShowKnownTypes",
		value: []int{1, 2, 3},
		configure: func(vw *CommonValueWriter) {
			vw.TypeAndValue.ShowKnownTypes = true
		},
	},
	{
		name:  "SliceNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseSlice.WriteValue}
		},
	},
	{
		name:  "String",
		value: "test",
	},
	{
		name:  "StringEmpty",
		value: "",
	},
	{
		name:  "StringUnquoted",
		value: "aaa\nbbb",
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseString.Quote = false
		},
	},
	{
		name:  "StringTruncated",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseString.MaxLen = 2
		},
	},
	{
		name:  "StringNot",
		value: 123,
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseString.WriteValue}
		},
	},
	{
		name: "Struct",
		value: testStruct{
			Foo:        123,
			Bar:        123.456,
			unexported: 123,
		},
	},
	{
		name:  "StructEmpty",
		value: struct{}{},
	},
	{
		name: "StructUnexportedDisabled",
		value: testStruct{
			Foo:        123,
			Bar:        123.456,
			unexported: 123,
		},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseStruct.Unexported = false
		},
	},
	{
		name:  "StructNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseStruct.WriteValue}
		},
	},
	{
		name: "UnsafePointer",
		value: func() unsafe.Pointer {
			i := 123
			return unsafe.Pointer(&i)
		}(),
		ignoreResult: true,
	},
	{
		name: "UnsafePointerNil",
		value: func() unsafe.Pointer {
			return unsafe.Pointer(nil)
		}(),
	},
	{
		name:  "UnsafePointerNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Kind.BaseUnsafePointer.WriteValue}
		},
	},
	{
		name:  "FilterMatch",
		value: &testStringer{s: "test"},
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter().WriteValue, func(v reflect.Value) bool {
				return v.Type() == reflect.TypeFor[*testStringer]()
			}).WriteValue}
		},
	},
	{
		name:  "FilterNoMatch",
		value: &testStringer{s: "test"},
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter().WriteValue, func(v reflect.Value) bool {
				return v.Type() != reflect.TypeFor[*testStringer]()
			}).WriteValue}
		},
	},
	{
		name:         "PanicString",
		value:        "test",
		panicRecover: true,
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st State, v reflect.Value) bool {
				panic("string")
			}}
		},
	},
	{
		name:         "PanicError",
		value:        "test",
		panicRecover: true,
		configure: func(vw *CommonValueWriter) {
			err := errors.New("error")
			vw.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st State, v reflect.Value) bool {
				panic(err)
			}}
		},
	},
	{
		name:         "PanicOther",
		value:        "test",
		panicRecover: true,
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st State, v reflect.Value) bool {
				panic(123)
			}}
		},
	},
	{
		name:         "PanicNot",
		value:        "test",
		panicRecover: true,
	},
	{
		name:  "UnwrapInterface",
		value: [1]any{123},
	},
	{
		name:  "UnwrapInterfaceNil",
		value: [1]any{},
	},
	{
		name:  "UnwrapInterfaceDisabled",
		value: [1]any{123},
		configure: func(vw *CommonValueWriter) {
			vw.UnwrapInterface = nil
		},
	},
	{
		name:  "UnwrapInterfaceDisabledNil",
		value: [1]any{},
		configure: func(vw *CommonValueWriter) {
			vw.UnwrapInterface = nil
		},
	},
	{
		name: "RecursionPointer",
		value: func() any {
			var v any
			v = &v
			return v
		}(),
	},
	{
		name: "RecursionSlice",
		value: func() []any {
			v := make([]any, 1)
			v[0] = v
			return v
		}(),
	},
	{
		name: "RecursionMap",
		value: func() map[int]any {
			v := make(map[int]any)
			v[0] = v
			return v
		}(),
	},
	{
		name:  "RecursionDisabled",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.Recursion = nil
		},
	},
	{
		name: "MaxDepth",
		value: func() any {
			var v1 any
			v2 := &v1
			v3 := &v2
			v4 := &v3
			return v4
		}(),
		configure: func(vw *CommonValueWriter) {
			vw.MaxDepth.Max = 2
		},
	},
	{
		name: "MaxDepthDisabled",
		value: func() any {
			var v1 any
			v2 := &v1
			v3 := &v2
			v4 := &v3
			return v4
		}(),
		configure: func(vw *CommonValueWriter) {
			vw.MaxDepth = nil
		},
	},
	{
		name: "TypeFullName",
		value: testStruct{
			Foo:        123,
			Bar:        123.456,
			unexported: 123,
		},
		configure: func(vw *CommonValueWriter) {
			vw.Type.Stringer = reflectutil.TypeFullName
		},
	},
	{
		name:  "TypeDisabled",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.Type = nil
		},
	},
	{
		name:  "TypeAndValueDisabled",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.TypeAndValue = nil
		},
	},
	{
		name:  "ReflectValue",
		value: reflect.ValueOf(123),
	},
	{
		name:  "ReflectValueNil",
		value: reflect.ValueOf(nil),
	},
	{
		name: "ReflectValueUnexported",
		value: testUnexported{
			v: reflect.ValueOf(123),
		},
	},
	{
		name:  "ReflectValueDisabled",
		value: reflect.ValueOf(123),
		configure: func(vw *CommonValueWriter) {
			vw.ReflectValue = nil
		},
		ignoreResult: true,
	},
	{
		name:  "ReflectValueNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.ReflectValue.WriteValue}
		},
	},
	{
		name:  "Error",
		value: &testError{},
	},
	{
		name:  "ErrorNil",
		value: (*testError)(nil),
	},
	{
		name:  "ErrorUnexported",
		value: testUnexported{v: &testError{}},
		configure: func(vw *CommonValueWriter) {
			vw.CanInterface = nil
		},
	},
	{
		name:  "ErrorUnexportedCanInterface",
		value: testUnexported{v: &testError{}},
	},
	{
		name:  "ErrorDisabled",
		value: &testError{},
		configure: func(vw *CommonValueWriter) {
			vw.Error = nil
		},
	},
	{
		name:  "ErrorNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Error.WriteValue}
		},
	},
	{
		name:  "BytesHexDump",
		value: bytes.Repeat([]byte("test"), 100),
	},
	{
		name:  "BytesHexDumpNil",
		value: []byte(nil),
	},
	{
		name:  "BytesHexDumpTruncated",
		value: []byte("test"),
		configure: func(vw *CommonValueWriter) {
			vw.BytesHexDump.MaxLen = 2
		},
	},
	{
		name:  "BytesHexDumpShowAddr",
		value: []byte("test"),
		configure: func(vw *CommonValueWriter) {
			vw.BytesHexDump.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name:  "BytesHexDumpDisabled",
		value: []byte("test"),
		configure: func(vw *CommonValueWriter) {
			vw.BytesHexDump = nil
		},
	},
	{
		name:  "BytesHexDumpNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.BytesHexDump.WriteValue}
		},
	},
	{
		name:  "BytesableHexDump",
		value: &testBytesable{b: bytes.Repeat([]byte("test"), 100)},
	},
	{
		name:  "BytesableHexDumpNil",
		value: (*testBytesable)(nil),
	},
	{
		name:  "BytesableHexDumpNilBytes",
		value: &testBytesable{},
	},
	{
		name:  "BytesableHexDumpTruncated",
		value: &testBytesable{b: []byte("test")},
		configure: func(vw *CommonValueWriter) {
			vw.BytesableHexDump.MaxLen = 2
		},
	},
	{
		name: "BytesableHexDumpUnexported",
		value: testUnexported{
			v: &testBytesable{},
		},
		configure: func(vw *CommonValueWriter) {
			vw.CanInterface = nil
		},
	},
	{
		name: "BytesableHexDumpUnexportedCanInterface",
		value: testUnexported{
			v: &testBytesable{},
		},
	},
	{
		name:  "BytesableHexDumpReflectValue",
		value: reflect.ValueOf(123),
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{
				vw.BytesableHexDump.WriteValue,
				vw.ReflectValue.WriteValue,
			}
		},
	},
	{
		name:  "BytesableHexDumpShowAddr",
		value: &testBytesable{b: []byte("test")},
		configure: func(vw *CommonValueWriter) {
			vw.BytesableHexDump.ShowAddr = true
		},
		ignoreResult: true,
	},
	{
		name:  "BytesableHexDumpDisabled",
		value: &testBytesable{b: []byte("test")},
		configure: func(vw *CommonValueWriter) {
			vw.BytesableHexDump = nil
		},
	},
	{
		name:  "BytesableHexDumpNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.BytesableHexDump.WriteValue}
		},
	},
	{
		name:  "Stringer",
		value: &testStringer{s: "test"},
	},
	{
		name:  "StringerNil",
		value: (*testStringer)(nil),
	},
	{
		name:  "StringerTruncated",
		value: &testStringer{s: "test"},
		configure: func(vw *CommonValueWriter) {
			vw.Stringer.MaxLen = 2
		},
	},
	{
		name:  "StringerUnexported",
		value: testUnexported{v: &testStringer{}},
		configure: func(vw *CommonValueWriter) {
			vw.CanInterface = nil
		},
	},
	{
		name:  "StringerUnexportedCanInterface",
		value: testUnexported{v: &testStringer{}},
	},
	{
		name:  "StringerReflectValue",
		value: reflect.ValueOf(123),
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{
				vw.Stringer.WriteValue,
				vw.ReflectValue.WriteValue,
			}
		},
	},
	{
		name:  "StringerDisabled",
		value: &testStringer{s: "test"},
		configure: func(vw *CommonValueWriter) {
			vw.Stringer = nil
		},
	},
	{
		name:  "StringerNot",
		value: "test",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{vw.Stringer.WriteValue}
		},
	},
	{
		name:  "DefaultPrinter",
		value: DefaultPrinter,
	},
	{
		name: "CommonValueWriter",
		value: func() *CommonValueWriter {
			vw := NewCommonValueWriter()
			vw.ConfigureTest()
			vw.SetShowLen(false)
			vw.SetShowCap(false)
			return vw
		}(),
	},
}

func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.newPrinter()
			s := p.String(tc.value)
			if !tc.ignoreResult {
				assertauto.Equal(t, s, assertauto.Name("result"))
			}
			t.Log(s)
			if !tc.ignoreAllocs {
				assertauto.AllocsPerRun(t, 100, func() {
					t.Helper()
					p.Write(io.Discard, tc.value)
				}, assertauto.Name("allocs"))
			}
		})
	}
}

func TestPanicWriterError(t *testing.T) {
	p, vw := newTestPrinter()
	vw.PanicRecover = nil
	w := &testErrorWriter{}
	assert.Panics(t, func() {
		p.Write(w, "test")
	})
}

func Benchmark(b *testing.B) {
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			p := tc.newPrinter()
			for range b.N {
				p.Write(io.Discard, tc.value)
			}
		})
	}
}

func TestPrinterPanicNotHandled(t *testing.T) {
	c := NewConfig()
	vw := func(c *Config, w io.Writer, st State, v reflect.Value) bool {
		return false
	}
	p := NewPrinter(c, vw)
	assert.Panics(t, func() {
		p.Write(io.Discard, "test")
	})
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	Write(buf, "test")
	s := buf.String()
	assertauto.Equal(t, s, assertauto.Name("result"))
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		Write(io.Discard, "test")
	}, assertauto.Name("allocs"))
}

func TestString(t *testing.T) {
	s := String("test")
	assertauto.Equal(t, s, assertauto.Name("result"))
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		String("test")
	}, assertauto.Name("allocs"))
}

func TestFormatter(t *testing.T) {
	buf := new(bytes.Buffer)
	f := Formatter("test")
	_, err := fmt.Fprintf(buf, "%v", f)
	assert.NoError(t, err)
	s := buf.String()
	assertauto.Equal(t, s, assertauto.Name("result"))
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		_, err := fmt.Fprintf(io.Discard, "%v", f)
		assert.NoError(t, err)
	}, assertauto.Name("allocs"))
}

var testIdent = DefaultConfig.Indent

func TestIndentWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	testIndentWriter(t, buf, testIdent, 1)
	assert.Equal(t, buf.String(), "\taabb\n\tc\n\tc\n\tdd") //nolint:dupword // Test data.
	assert.AllocsPerRun(t, 100, func() {
		t.Helper()
		testIndentWriter(t, io.Discard, testIdent, 1)
	}, 0)
}

var testIndentWriterValues = []struct {
	b         []byte
	expectedN int
}{
	{
		b:         []byte("aa"),
		expectedN: 3,
	},
	{
		b:         []byte("bb\n"),
		expectedN: 3,
	},
	{
		b:         []byte("c\nc"),
		expectedN: 5,
	},
	{
		b:         []byte("\ndd"),
		expectedN: 4,
	},
}

func testIndentWriter(tb testing.TB, w io.Writer, indent string, level int) {
	tb.Helper()
	iw := NewIndentWriter(w, indent, level, false)
	for _, v := range testIndentWriterValues {
		n, err := iw.Write(v.b)
		assert.NoError(tb, err)
		assert.Equal(tb, n, v.expectedN)
	}
}

func TestIndentWriterErrorIndent(t *testing.T) {
	w := &testErrorWriter{}
	iw := NewIndentWriter(w, DefaultConfig.Indent, 1, false)
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func TestIndentWriterErrorWrite(t *testing.T) {
	w := &testErrorWriter{}
	iw := NewIndentWriter(w, DefaultConfig.Indent, 1, true)
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func BenchmarkIndentWriter(b *testing.B) {
	iw := NewIndentWriter(io.Discard, DefaultConfig.Indent, 1, false)
	for range b.N {
		for _, v := range testIndentWriterValues {
			_, _ = iw.Write(v.b)
		}
	}
}

var testWriteIndentLevels = []int{0, 1, 2, 5, 10, 20, 50, 100, 200, 500, 1000}

func TestWriteIndent(t *testing.T) {
	for _, level := range testWriteIndentLevels {
		t.Run(strconv.Itoa(level), func(t *testing.T) {
			buf := new(bytes.Buffer)
			WriteIndent(buf, testIdent, level)
			assert.Equal(t, buf.String(), strings.Repeat(testIdent, level))
			assert.AllocsPerRun(t, 100, func() {
				WriteIndent(io.Discard, testIdent, level)
			}, 0)
		})
	}
}

func BenchmarkWriteIndent(b *testing.B) {
	for _, level := range testWriteIndentLevels {
		b.Run(strconv.Itoa(level), func(b *testing.B) {
			for range b.N {
				WriteIndent(io.Discard, testIdent, level)
			}
		})
	}
}

type testStruct struct {
	Foo        int
	Bar        float64
	unexported int
}

type testComparableStruct struct {
	V int
}

type testUnexported struct {
	v any
}

type testBytesable struct {
	b []byte
}

func (b *testBytesable) Bytes() []byte {
	return b.b
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}

type testStringer struct {
	s string
}

func (sr *testStringer) String() string {
	return sr.s
}

type testErrorWriter struct{}

func (w *testErrorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("test")
}
