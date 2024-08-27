package pretty_test

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty"
)

var testIndent = DefaultConfig.Indent

func TestIndentWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	testIndentWriter(t, buf, testIndent, 1)
	assert.Equal(t, buf.String(), "\taabb\n\tc\n\tc\n\tdd") //nolint:dupword // Test data.
	assert.AllocsPerRun(t, 100, func() {
		t.Helper()
		testIndentWriter(t, io.Discard, testIndent, 1)
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

var (
	writeIndentTestCases = []struct {
		name   string
		indent string
	}{
		{
			name:   "Tab",
			indent: "\t",
		},
		{
			name:   "4Spaces",
			indent: "    ",
		},
	}
	testWriteIndentLevels = []int{0, 1, 2, 10, 100, 1000, 1001}
)

func TestWriteIndent(t *testing.T) {
	for _, tc := range writeIndentTestCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, level := range testWriteIndentLevels {
				t.Run(strconv.Itoa(level), func(t *testing.T) {
					buf := new(bytes.Buffer)
					WriteIndent(buf, tc.indent, level)
					assert.Equal(t, buf.String(), strings.Repeat(tc.indent, level))
					assert.AllocsPerRun(t, 100, func() {
						WriteIndent(io.Discard, tc.indent, level)
					}, 0)
				})
			}
		})
	}
}

func BenchmarkWriteIndent(b *testing.B) {
	for _, tc := range writeIndentTestCases {
		b.Run(tc.name, func(b *testing.B) {
			for _, level := range testWriteIndentLevels {
				b.Run(strconv.Itoa(level), func(b *testing.B) {
					for range b.N {
						WriteIndent(io.Discard, tc.indent, level)
					}
				})
			}
		})
	}
}

func BenchmarkWriteIndentParallel(b *testing.B) {
	for _, tc := range writeIndentTestCases {
		b.Run(tc.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					WriteIndent(io.Discard, tc.indent, 10)
				}
			})
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
