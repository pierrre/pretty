package pretty_test

import (
	"bytes"
	"errors"
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

type testCase struct {
	name         string
	value        any
	panicRecover bool
	configure    func(vw *CommonValueWriter)
	ignoreAllocs bool
}

func (tc testCase) newPrinter() *Printer {
	c := NewConfig()
	vw := NewCommonValueWriter()
	if !tc.panicRecover {
		vw.PanicRecover = nil
	}
	p := NewPrinterCommon(c, vw)
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
		value: "dummy",
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
				return vw.Kind.WriteValue(c, w, st, reflect.ValueOf(nil))
			}}
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
		name: "ChanShowLenOnly",
		value: func() chan int {
			c := make(chan int, 5)
			c <- 123
			return c
		}(),
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseChan.ShowCap = false
		},
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
		name:  "MapUnsorted",
		value: map[int]int{1: 2},
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
		name:  "SliceShowLenOnly",
		value: []int{1, 2, 3},
		configure: func(vw *CommonValueWriter) {
			vw.Kind.BaseSlice.ShowCap = false
		},
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
			var zero unsafe.Pointer
			return zero
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
			vw.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
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
			vw.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
				panic(err)
			}}
		},
	},
	{
		name:         "PanicOther",
		value:        "test",
		panicRecover: true,
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
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
		name:  "BytesHex",
		value: bytes.Repeat([]byte("test"), 100),
	},
	{
		name:  "BytesHexNil",
		value: []byte(nil),
	},
	{
		name:  "BytesHexTruncated",
		value: []byte("test"),
		configure: func(vw *CommonValueWriter) {
			vw.BytesHex.MaxLen = 2
		},
	},
	{
		name:  "BytesableHex",
		value: &testBytesable{b: bytes.Repeat([]byte("test"), 100)},
	},
	{
		name:  "BytesableHexNil",
		value: (*testBytesable)(nil),
	},
	{
		name:  "BytesableHexNilBytes",
		value: &testBytesable{},
	},
	{
		name:  "BytesableHexTruncated",
		value: &testBytesable{b: []byte("test")},
		configure: func(vw *CommonValueWriter) {
			vw.BytesableHex.MaxLen = 2
		},
	},
	{
		name: "BytesableHexUnexported",
		value: testUnexported{
			v: &testBytesable{},
		},
	},
	{
		name:  "BytesableHexReflectValue",
		value: reflect.ValueOf(123),
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{
				vw.BytesableHex.WriteValue,
				vw.Reflect.WriteValue,
			}
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
	},
	{
		name:  "StringerReflectValue",
		value: reflect.ValueOf(123),
		configure: func(vw *CommonValueWriter) {
			vw.ValueWriters = ValueWriters{
				vw.Stringer.WriteValue,
				vw.Reflect.WriteValue,
			}
		},
	},
}

func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.newPrinter()
			s := p.String(tc.value)
			assertauto.Equal(t, s, assertauto.Name("result"))
			t.Log(s)
			if !tc.ignoreAllocs {
				allocs := testing.AllocsPerRun(100, func() {
					t.Helper()
					p.Write(io.Discard, tc.value)
				})
				assertauto.Equal(t, allocs, assertauto.Name("allocs"))
			}
		})
	}
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

func TestIndentWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewConfig()
	st := &State{
		Indent: 1,
	}
	testIndentWriter(t, c, buf, st)
	assert.Equal(t, buf.String(), "\taabb\n\tc\n\tc\n\tdd") //nolint:dupword // Test data.
	assert.AllocsPerRun(t, 100, func() {
		t.Helper()
		testIndentWriter(t, c, io.Discard, st)
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

func testIndentWriter(tb testing.TB, c *Config, w io.Writer, st *State) {
	tb.Helper()
	iw := NewIndentWriter(c, w, st, false)
	for _, v := range testIndentWriterValues {
		n, err := iw.Write(v.b)
		assert.NoError(tb, err)
		assert.Equal(tb, n, v.expectedN)
	}
}

func TestIndentWriterErrorIndent(t *testing.T) {
	w := &testErrorWriter{}
	c := NewConfig()
	st := &State{
		Indent: 1,
	}
	iw := NewIndentWriter(c, w, st, false)
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func TestIndentWriterErrorWrite(t *testing.T) {
	w := &testErrorWriter{}
	c := NewConfig()
	st := &State{
		Indent: 1,
	}
	iw := NewIndentWriter(c, w, st, true)
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func BenchmarkIndentWriter(b *testing.B) {
	c := NewConfig()
	st := &State{
		Indent: 1,
	}
	iw := NewIndentWriter(c, io.Discard, st, false)
	for range b.N {
		for _, v := range testIndentWriterValues {
			_, _ = iw.Write(v.b)
		}
	}
}

var (
	testWriteIndentString = "\t"
	testWriteIndentCounts = []int{0, 1, 2, 5, 10, 20, 50, 100, 200, 500, 1000}
)

func TestWriteIndent(t *testing.T) {
	for _, count := range testWriteIndentCounts {
		t.Run(strconv.Itoa(count), func(t *testing.T) {
			buf := new(bytes.Buffer)
			WriteIndent(buf, testWriteIndentString, count)
			assert.Equal(t, buf.String(), strings.Repeat(testWriteIndentString, count))
			assert.AllocsPerRun(t, 100, func() {
				WriteIndent(io.Discard, testWriteIndentString, count)
			}, 0)
		})
	}
}

func BenchmarkWriteIndent(b *testing.B) {
	for _, count := range []int{0, 1, 2, 5, 10, 20, 50, 100, 200, 500, 1000} {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			for range b.N {
				WriteIndent(io.Discard, testWriteIndentString, count)
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
