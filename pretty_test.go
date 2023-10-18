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
		_, _ = io.WriteString(w, "example")
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
		name:  "Bytes",
		value: bytes.Repeat([]byte("test"), 100),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewBytesValueWriter(0)}
		},
	},
	{
		name:  "BytesNil",
		value: []byte(nil),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewBytesValueWriter(0)}
		},
	},
	{
		name:  "BytesTruncated",
		value: []byte("test"),
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewBytesValueWriter(2)}
		},
	},
	{
		name:  "Byteser",
		value: &testByteser{b: bytes.Repeat([]byte("test"), 100)},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserValueWriter(0)}
		},
	},
	{
		name:  "ByteserNil",
		value: &testByteser{},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserValueWriter(0)}
		},
	},
	{
		name:  "ByteserTruncated",
		value: &testByteser{b: []byte("test")},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserValueWriter(2)}
		},
	},
	{
		name: "ByteserUnexported",
		value: testUnexported{
			v: &testByteser{},
		},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewByteserValueWriter(0)}
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
		name:  "FilterMatch",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter(0), func(v reflect.Value) bool {
				return v.Type() == reflect.TypeOf(&testStringer{})
			})}
		},
	},
	{
		name:  "FilterNoMatch",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter(0), func(v reflect.Value) bool {
				return v.Type() != reflect.TypeOf(&testStringer{})
			})}
		},
	},
	{
		name:  "Default",
		value: &testStringer{s: "test"},
		configure: func(c *Config) {
			c.ValueWriters = []ValueWriter{
				NewDefaultValueWriter(),
				NewStringerValueWriter(0),
			}
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
			for i := 0; i < b.N; i++ {
				c.Write(io.Discard, tc.value)
			}
		})
	}
}

var writeIndentValues = []int{0, 1, 2, 5, 10, 20, 50, 100, 200, 500, 1000}

func TestConfigWriteIndent(t *testing.T) {
	c := newTestConfig()
	for _, indent := range writeIndentValues {
		t.Run(strconv.Itoa(indent), func(t *testing.T) {
			st := &State{
				Indent: indent,
			}
			buf := new(bytes.Buffer)
			c.WriteIndent(buf, st)
			assert.Equal(t, buf.String(), strings.Repeat(c.Indent, indent))
			allocs := testing.AllocsPerRun(100, func() {
				t.Helper()
				c.WriteIndent(io.Discard, st)
			})
			assert.Equal(t, allocs, 0)
		})
	}
}

func BenchmarkConfigWriteIndent(b *testing.B) {
	c := newTestConfig()
	for _, indent := range []int{0, 1, 2, 5, 10, 20, 50, 100, 200, 500, 1000} {
		b.Run(strconv.Itoa(indent), func(b *testing.B) {
			st := &State{
				Indent: indent,
			}
			for i := 0; i < b.N; i++ {
				c.WriteIndent(io.Discard, st)
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

func TestIndentWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	c := newTestConfig()
	st := &State{
		Indent: 1,
	}
	testIndentWriter(t, c, buf, st)
	assertauto.Equal(t, buf.String())
	allocs := testing.AllocsPerRun(100, func() {
		t.Helper()
		testIndentWriter(t, c, io.Discard, st)
	})
	assertauto.Equal(t, allocs)
}

var testIndentWriterValues = [][]byte{
	[]byte("aa"),
	[]byte("bb\n"),
	[]byte("c\nc"),
	[]byte("\ndd"),
}

func testIndentWriter(tb testing.TB, c *Config, w io.Writer, st *State) {
	tb.Helper()
	iw := GetIndentWriter(w, c, st, false)
	for _, b := range testIndentWriterValues {
		n, err := iw.Write(b)
		assert.NoError(tb, err)
		assert.Equal(tb, n, len(b))
	}
	iw.Release()
}

func TestIndentWriterError(t *testing.T) {
	w := &testErrorWriter{}
	c := newTestConfig()
	st := &State{
		Indent: 1,
	}
	iw := GetIndentWriter(w, c, st, false)
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
	iw.Release()
}

func BenchmarkIndentWriter(b *testing.B) {
	c := newTestConfig()
	st := &State{
		Indent: 1,
	}
	iw := GetIndentWriter(io.Discard, c, st, false)
	for i := 0; i < b.N; i++ {
		for _, b := range testIndentWriterValues {
			_, _ = iw.Write(b)
		}
	}
	iw.Release()
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
