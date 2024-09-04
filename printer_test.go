package pretty_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	. "github.com/pierrre/pretty"
)

func newTestPrinter() (*Printer, *CommonValueWriter) {
	c := NewConfig()
	vw := NewCommonValueWriter()
	vw.ConfigureTest()
	p := NewPrinter(c, vw.WriteValue)
	return p, vw
}

type testCase struct {
	name            string
	value           any
	panicRecover    bool
	configure       func(vw *CommonValueWriter)
	options         []Option
	ignoreResult    bool
	ignoreAllocs    bool
	ignoreBenchmark bool
}

func (tc *testCase) newPrinter() *Printer {
	p, vw := newTestPrinter()
	if !tc.panicRecover {
		vw.PanicRecover = nil
	}
	if tc.configure != nil {
		tc.configure(vw)
	}
	return p
}

var testCases = []*testCase{
	{
		name:  "Options",
		value: 123,
		options: []Option{func(st *State) {
			st.KnownType = true
		}},
	},
}

func addTestCases(tcs []*testCase) {
	testCases = append(testCases, tcs...)
}

func addTestCasesPrefix(prefix string, tcs []*testCase) {
	for _, tc := range tcs {
		tc.name = prefix + "/" + tc.name
	}
	addTestCases(tcs)
}

func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.newPrinter()
			s := p.String(tc.value, tc.options...)
			if !tc.ignoreResult {
				assertauto.Equal(t, s, assertauto.Name("result"))
			}
			t.Log(s)
			if !tc.ignoreAllocs {
				allocs, _ := assertauto.AllocsPerRun(t, 100, func() {
					t.Helper()
					p.Write(io.Discard, tc.value)
				}, assertauto.Name("allocs"))
				t.Logf("allocs: %g", allocs)
			}
		})
	}
}

func TestPrinterPanicWriterError(t *testing.T) {
	p, vw := newTestPrinter()
	vw.PanicRecover = nil
	w := &testErrorWriter{}
	assert.Panics(t, func() {
		p.Write(w, "test")
	})
}

func TestPrinterPanicNotHandled(t *testing.T) {
	c := NewConfig()
	vw := func(st *State, v reflect.Value) bool {
		return false
	}
	p := NewPrinter(c, vw)
	assert.Panics(t, func() {
		p.Write(io.Discard, "test")
	})
}

func Benchmark(b *testing.B) {
	for _, tc := range testCases {
		if tc.ignoreBenchmark {
			continue
		}
		b.Run(tc.name, func(b *testing.B) {
			p := tc.newPrinter()
			for range b.N {
				p.Write(io.Discard, tc.value, tc.options...)
			}
		})
	}
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

type testErrorWriter struct{}

func (w *testErrorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("test")
}
