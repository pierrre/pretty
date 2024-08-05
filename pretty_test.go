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
	. "github.com/pierrre/pretty"
)

func init() {
	DefaultConfig.MapSortKeys = true // Makes tests more deterministic.
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
	// Output: (pretty_test.exampleStruct) {
	// 	Int: (int) 123,
	// 	Float: (float64) 123.456,
	// 	String: (string) (len=4) "test",
	// 	Map: (map[string]int) (len=2) {
	// 		(string) (len=3) "bar": (int) 2,
	// 		(string) (len=3) "foo": (int) 1,
	// 	},
	// 	Slice: ([]int) (len=3 cap=3) {
	// 		(int) 1,
	// 		(int) 2,
	// 		(int) 3,
	// 	},
	// }
}

func ExampleString() {
	s := String("test")
	fmt.Println(s)
	// Output: (string) (len=4) "test"
}

func ExampleWrite() {
	buf := new(bytes.Buffer)
	Write(buf, "test")
	s := buf.String()
	fmt.Println(s)
	// Output: (string) (len=4) "test"
}

func ExampleFormatter() {
	f := Formatter("test")
	s := fmt.Sprintf("%v", f)
	fmt.Println(s)
	// Output: (string) (len=4) "test"
}

func ExampleValueWriter() {
	vw := func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		WriteString(w, "example")
		return true
	}
	c := NewConfig()
	c.ValueWriters = []ValueWriter{vw}
	s := c.String("test")
	fmt.Println(s)
	// Output: (string) example
}

func newTestConfig() *Config {
	c := NewConfig()
	c.ValueWriters = nil
	return c
}

type testCase struct {
	name      string
	value     any
	configure func(c *Config)
}

func (tc testCase) newConfig() *Config {
	c := newTestConfig()
	if tc.configure != nil {
		tc.configure(c)
	}
	return c
}

var testCases = []testCase{
	{
		name:  "Nil",
		value: nil,
	},
	{
		name:  "PanicString",
		value: "test",
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
				panic("string")
			}}
		},
	},
	{
		name:  "PanicError",
		value: "test",
		configure: func(c *Config) {
			err := errors.New("error")
			c.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
				panic(err)
			}}
		},
	},
	{
		name:  "PanicOther",
		value: "test",
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
				panic(123)
			}}
		},
	},
	{
		name:  "Bool",
		value: true,
	},
	{
		name:  "Int",
		value: 123,
	},
	{
		name:  "Uint",
		value: uint(123),
	},
	{
		name:  "Uintptr",
		value: uintptr(123),
	},
	{
		name:  "Float",
		value: 123.456,
	},
	{
		name:  "Complex",
		value: 123.456 + 789.123i,
	},
	{
		name:  "String",
		value: "test",
	},
	{
		name:  "StringTruncated",
		value: "test",
		configure: func(c *Config) {
			c.StringMaxLen = 2
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
		name:  "Func",
		value: NewConfig,
	},
	{
		name:  "FuncNil",
		value: (func())(nil),
	},
	{
		name:  "Interface",
		value: [1]any{123},
	},
	{
		name: "Pointer",
		value: func() *int {
			i := 123
			return &i
		}(),
	},
	{
		name: "UnsafePointer",
		value: func() unsafe.Pointer {
			var zero unsafe.Pointer
			return zero
		}(),
	},
	{
		name:  "Array",
		value: [3]int{1, 2, 3},
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
		name:  "SliceEmpty",
		value: []int{},
	},
	{
		name:  "SliceTruncated",
		value: []int{1, 2, 3},
		configure: func(c *Config) {
			c.SliceMaxLen = 2
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
	},
	{
		name:  "MapUnsortedUnexported",
		value: testUnexported{v: map[int]int{1: 2}},
	},
	{
		name:  "MapSortedBool",
		value: map[bool]int{false: 1, true: 2},
		configure: func(c *Config) {
			c.MapSortKeys = true
		},
	},
	{
		name:  "MapSortedInt",
		value: map[int]int{1: 2, 3: 4, 5: 6},
		configure: func(c *Config) {
			c.MapSortKeys = true
		},
	},
	{
		name:  "MapSortedUint",
		value: map[uint]int{1: 2, 3: 4, 5: 6},
		configure: func(c *Config) {
			c.MapSortKeys = true
		},
	},
	{
		name:  "MapSortedFloat",
		value: map[float64]int{1: 2, 3: 4, 5: 6},
		configure: func(c *Config) {
			c.MapSortKeys = true
		},
	},
	{
		name:  "MapSortedString",
		value: map[string]int{"a": 1, "b": 2, "c": 3},
		configure: func(c *Config) {
			c.MapSortKeys = true
		},
	},
	{
		name:  "MapSortedDefault",
		value: map[testComparableStruct]int{{V: 1}: 2, {V: 3}: 4, {V: 5}: 6},
		configure: func(c *Config) {
			c.MapSortKeys = true
		},
	},
	{
		name:  "MapSortedTruncated",
		value: map[int]int{1: 2, 3: 4, 5: 6},
		configure: func(c *Config) {
			c.MapSortKeys = true
			c.MapMaxLen = 2
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
		name: "StructUnexportedDisabled",
		value: testStruct{
			Foo:        123,
			Bar:        123.456,
			unexported: 123,
		},
		configure: func(c *Config) {
			c.StructUnexported = false
		},
	},
	{
		name: "TypeFullName",
		value: testStruct{
			Foo:        123,
			Bar:        123.456,
			unexported: 123,
		},
		configure: func(c *Config) {
			c.TypeFullName = true
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
		configure: func(c *Config) {
			c.MaxDepth = 2
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
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewReflectValueValueWriter()}
		},
	},
	{
		name: "ReflectValueUnexported",
		value: testUnexported{
			v: reflect.ValueOf(123),
		},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewReflectValueValueWriter()}
		},
	},
	{
		name:  "BytesHex",
		value: bytes.Repeat([]byte("test"), 100),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewBytesHexValueWriter(0)}
		},
	},
	{
		name:  "BytesHexNil",
		value: []byte(nil),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewBytesHexValueWriter(0)}
		},
	},
	{
		name:  "BytesHexTruncated",
		value: []byte("test"),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewBytesHexValueWriter(2)}
		},
	},
	{
		name:  "ByteserHex",
		value: &testByteser{b: bytes.Repeat([]byte("test"), 100)},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserHexValueWriter(0)}
		},
	},
	{
		name:  "ByteserHexNil",
		value: (*testByteser)(nil),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserHexValueWriter(0)}
		},
	},
	{
		name:  "ByteserHexNilBytes",
		value: &testByteser{},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserHexValueWriter(0)}
		},
	},
	{
		name:  "ByteserHexTruncated",
		value: &testByteser{b: []byte("test")},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserHexValueWriter(2)}
		},
	},
	{
		name: "ByteserHexUnexported",
		value: testUnexported{
			v: &testByteser{},
		},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserHexValueWriter(0)}
		},
	},
	{
		name:  "ByteserHexReflectValue",
		value: reflect.ValueOf(123),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{
				NewByteserHexValueWriter(0),
				NewReflectValueValueWriter(),
			}
		},
	},
	{
		name:  "Error",
		value: &testError{},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewErrorValueWriter()}
		},
	},
	{
		name:  "ErrorNil",
		value: (*testError)(nil),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewErrorValueWriter()}
		},
	},
	{
		name:  "ErrorUnexported",
		value: testUnexported{v: &testError{}},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewErrorValueWriter()}
		},
	},
	{
		name:  "Stringer",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewStringerValueWriter(0)}
		},
	},
	{
		name:  "StringerNil",
		value: (*testStringer)(nil),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewStringerValueWriter(0)}
		},
	},
	{
		name:  "StringerTruncated",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewStringerValueWriter(2)}
		},
	},
	{
		name:  "StringerUnexported",
		value: testUnexported{v: &testStringer{}},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewStringerValueWriter(0)}
		},
	},
	{
		name:  "StringerReflectValue",
		value: reflect.ValueOf(123),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{
				NewStringerValueWriter(0),
				NewReflectValueValueWriter(),
			}
		},
	},
	{
		name:  "FilterMatch",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter(0), func(v reflect.Value) bool {
				return v.Type() == reflect.TypeFor[*testStringer]()
			})}
		},
	},
	{
		name:  "FilterNoMatch",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter(0), func(v reflect.Value) bool {
				return v.Type() != reflect.TypeFor[*testStringer]()
			})}
		},
	},
}

func TestConfig(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.newConfig()
			s := c.String(tc.value)
			assertauto.Equal(t, s)
			t.Log(s)
			allocs := testing.AllocsPerRun(100, func() {
				t.Helper()
				c.Write(io.Discard, tc.value)
			})
			assertauto.Equal(t, allocs)
		})
	}
}

func BenchmarkConfig(b *testing.B) {
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			c := tc.newConfig()
			for range b.N {
				c.Write(io.Discard, tc.value)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	Write(buf, "test")
	s := buf.String()
	assertauto.Equal(t, s)
	allocs := testing.AllocsPerRun(100, func() {
		t.Helper()
		Write(io.Discard, "test")
	})
	assertauto.Equal(t, allocs)
}

func TestString(t *testing.T) {
	s := String("test")
	assertauto.Equal(t, s)
	allocs := testing.AllocsPerRun(100, func() {
		t.Helper()
		String("test")
	})
	assertauto.Equal(t, allocs)
}

func TestFormatter(t *testing.T) {
	buf := new(bytes.Buffer)
	f := Formatter("test")
	_, err := fmt.Fprintf(buf, "%v", f)
	assert.NoError(t, err)
	s := buf.String()
	assertauto.Equal(t, s)
	allocs := testing.AllocsPerRun(100, func() {
		t.Helper()
		_, err := fmt.Fprintf(io.Discard, "%v", f)
		assert.NoError(t, err)
	})
	assertauto.Equal(t, allocs)
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

func TestIndentWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	c := newTestConfig()
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
	iw := GetIndentWriter(w, c, st, false)
	defer iw.Release()
	for _, v := range testIndentWriterValues {
		n, err := iw.Write(v.b)
		assert.NoError(tb, err)
		assert.Equal(tb, n, v.expectedN)
	}
}

func TestIndentWriterErrorIndent(t *testing.T) {
	w := &testErrorWriter{}
	c := newTestConfig()
	st := &State{
		Indent: 1,
	}
	iw := GetIndentWriter(w, c, st, false)
	defer iw.Release()
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func TestIndentWriterErrorWrite(t *testing.T) {
	w := &testErrorWriter{}
	c := newTestConfig()
	st := &State{
		Indent: 1,
	}
	iw := GetIndentWriter(w, c, st, true)
	defer iw.Release()
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func BenchmarkIndentWriter(b *testing.B) {
	c := newTestConfig()
	st := &State{
		Indent: 1,
	}
	iw := GetIndentWriter(io.Discard, c, st, false)
	defer iw.Release()
	for range b.N {
		for _, v := range testIndentWriterValues {
			_, _ = iw.Write(v.b)
		}
	}
}

func TestNoErrorPanic(t *testing.T) {
	assert.Panics(t, func() {
		WriteBytes(&testErrorWriter{}, []byte("test"))
	})
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

type testByteser struct {
	b []byte
}

func (b *testByteser) Bytes() []byte {
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
