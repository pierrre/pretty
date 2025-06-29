package pretty_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Printer", []*prettytest.Case{
		{
			Name:         "Default",
			Value:        DefaultPrinter,
			IgnoreResult: true,
		},
	})
}

func Test(t *testing.T) {
	prettytest.Test(t)
}

func Benchmark(b *testing.B) {
	prettytest.Benchmark(b)
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
