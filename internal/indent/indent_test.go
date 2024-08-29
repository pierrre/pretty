package indent_test

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty/internal/indent"
)

var (
	writeTestCases = []struct {
		name   string
		string string
	}{
		{
			name:   "Tab",
			string: "\t",
		},
		{
			name:   "4Spaces",
			string: "    ",
		},
	}
	testWriteLevels = []int{0, 1, 2, 10, 100, 1000, 1001}
)

func TestWrite(t *testing.T) {
	for _, tc := range writeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, level := range testWriteLevels {
				t.Run(strconv.Itoa(level), func(t *testing.T) {
					buf := new(bytes.Buffer)
					MustWrite(buf, tc.string, level)
					assert.Equal(t, buf.String(), strings.Repeat(tc.string, level))
					assert.AllocsPerRun(t, 100, func() {
						MustWrite(io.Discard, tc.string, level)
					}, 0)
				})
			}
		})
	}
}

func BenchmarkWrite(b *testing.B) {
	for _, tc := range writeTestCases {
		b.Run(tc.name, func(b *testing.B) {
			for _, level := range testWriteLevels {
				b.Run(strconv.Itoa(level), func(b *testing.B) {
					for range b.N {
						MustWrite(io.Discard, tc.string, level)
					}
				})
			}
		})
	}
}

func BenchmarkWriteParallel(b *testing.B) {
	for _, tc := range writeTestCases {
		b.Run(tc.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					MustWrite(io.Discard, tc.string, 10)
				}
			})
		})
	}
}

func TestWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	testWriter(t, buf, Default, 1)
	assert.Equal(t, buf.String(), "\taabb\n\tc\n\tc\n\tdd") //nolint:dupword // Test data.
	assert.AllocsPerRun(t, 100, func() {
		t.Helper()
		testWriter(t, io.Discard, Default, 1)
	}, 0)
}

var testWriterValues = []struct {
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

func testWriter(tb testing.TB, w io.Writer, str string, level int) {
	tb.Helper()
	iw := NewWriter(w, str, level, false)
	defer iw.Release()
	for _, v := range testWriterValues {
		n, err := iw.Write(v.b)
		assert.NoError(tb, err)
		assert.Equal(tb, n, v.expectedN)
	}
}

func TestWriterErrorUnindented(t *testing.T) {
	w := &testErrorWriter{}
	iw := NewWriter(w, Default, 1, false)
	defer iw.Release()
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func TestWriterErrorIndented(t *testing.T) {
	w := &testErrorWriter{}
	iw := NewWriter(w, Default, 1, true)
	n, err := iw.Write([]byte("test"))
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func BenchmarkWriter(b *testing.B) {
	iw := NewWriter(io.Discard, Default, 1, false)
	defer iw.Release()
	for range b.N {
		for _, v := range testWriterValues {
			_, _ = iw.Write(v.b)
		}
	}
}

type testErrorWriter struct{}

func (w *testErrorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("test")
}
