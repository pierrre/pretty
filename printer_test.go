package pretty_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	. "github.com/pierrre/pretty"
)

type testCase struct {
	name             string
	value            any
	configurePrinter func(p *Printer)
	configureWriter  func(vw *CommonValueWriter)
	ignoreResult     bool
	ignoreAllocs     bool
	ignoreBenchmark  bool
}

func (tc *testCase) newPrinter() *Printer {
	vw := NewCommonValueWriter()
	vw.ConfigureTest(true)
	if tc.configureWriter != nil {
		tc.configureWriter(vw)
	}
	p := NewPrinter(vw)
	if tc.configurePrinter != nil {
		tc.configurePrinter(p)
	}
	return p
}

var testCases []*testCase

func addTestCases(tcs []*testCase) {
	testCases = append(testCases, tcs...)
}

func addTestCasesPrefix(prefix string, tcs []*testCase) {
	for _, tc := range tcs {
		tc.name = prefix + "/" + tc.name
	}
	addTestCases(tcs)
}

func init() {
	addTestCasesPrefix("Printer", []*testCase{
		{
			name:         "Default",
			value:        DefaultPrinter,
			ignoreResult: true,
		},
	})
}

func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.newPrinter()
			if !tc.ignoreResult {
				assertauto.Equal(t, tc.value, assertauto.ValueStringer(p.String))
			}
			s := p.String(tc.value)
			t.Log(s)
			if !tc.ignoreAllocs {
				allocs, _ := assertauto.AllocsPerRun(t, 100, func() {
					t.Helper()
					p.Write(io.Discard, tc.value)
				})
				t.Logf("allocs: %g", allocs)
			}
		})
	}
}

func Benchmark(b *testing.B) {
	for _, tc := range testCases {
		if tc.ignoreBenchmark {
			continue
		}
		b.Run(tc.name, func(b *testing.B) {
			p := tc.newPrinter()
			for b.Loop() {
				p.Write(io.Discard, tc.value)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	Write(buf, "test")
	s := buf.String()
	assertauto.Equal(t, s)
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		Write(io.Discard, "test")
	})
}

func TestWriteErr(t *testing.T) {
	buf := new(bytes.Buffer)
	err := WriteErr(buf, "test")
	assert.NoError(t, err)
	s := buf.String()
	assertauto.Equal(t, s)
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		Write(io.Discard, "test")
	})
}

func TestWriteErrError(t *testing.T) {
	w := &testPanicWriter{}
	err := WriteErr(w, "test")
	assert.Error(t, err)
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		_ = WriteErr(w, "test")
	})
}

type testPanicWriter struct{}

func (w *testPanicWriter) Write(p []byte) (int, error) {
	panic("test")
}

func TestString(t *testing.T) {
	s := String("test")
	assertauto.Equal(t, s)
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		String("test")
	})
}

func TestFormatter(t *testing.T) {
	buf := new(bytes.Buffer)
	f := Formatter("test")
	_, err := fmt.Fprintf(buf, "%v", f)
	assert.NoError(t, err)
	s := buf.String()
	assertauto.Equal(t, s)
	assertauto.AllocsPerRun(t, 100, func() {
		t.Helper()
		_, err := fmt.Fprintf(io.Discard, "%v", f)
		assert.NoError(t, err)
	})
}

type testUnexported struct {
	v any
}
