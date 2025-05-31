package write_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty/internal/write"
)

func TestMust(t *testing.T) {
	assert.NotPanics(t, func() {
		Must(0, nil)
	})
	assert.Panics(t, func() {
		Must(0, errors.New("test"))
	})
}

func TestString(t *testing.T) {
	buf := new(bytes.Buffer)
	n, err := String(buf, "test")
	assert.NoError(t, err)
	assert.Equal(t, n, 4)
	assert.Equal(t, buf.String(), "test")
}

func TestStringError(t *testing.T) {
	n, err := String(&testErrorWriter{}, "test")
	assert.Error(t, err)
	assert.Equal(t, n, 0)
}

func TestMustString(t *testing.T) {
	buf := new(bytes.Buffer)
	MustString(buf, "test")
	assert.Equal(t, buf.String(), "test")
}

func TestMustStringPanic(t *testing.T) {
	assert.Panics(t, func() {
		MustString(&testErrorWriter{}, "test")
	})
}

type testErrorWriter struct{}

func (w *testErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("error")
}
