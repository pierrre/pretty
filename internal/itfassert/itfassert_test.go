package itfassert_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty/internal/itfassert"
)

func Test(t *testing.T) {
	v := reflect.ValueOf(new(bytes.Buffer))
	w, ok := Assert[io.Writer](v)
	assert.True(t, ok)
	_, err := w.Write([]byte("test"))
	assert.NoError(t, err)
}

func TestFailUnexported(t *testing.T) {
	type unexportedStruct struct {
		unexported *bytes.Buffer
	}
	v := reflect.ValueOf(unexportedStruct{unexported: new(bytes.Buffer)}).FieldByName("unexported")
	_, ok := Assert[io.Writer](v)
	assert.False(t, ok)
}

func TestFailPointerNil(t *testing.T) {
	v := reflect.ValueOf((*bytes.Buffer)(nil))
	_, ok := Assert[io.Writer](v)
	assert.False(t, ok)
}

func TestFailAssertion(t *testing.T) {
	v := reflect.ValueOf("test")
	_, ok := Assert[io.Writer](v)
	assert.False(t, ok)
}

func TestFailAssertionInterface(t *testing.T) {
	v := reflect.ValueOf([1]io.Writer{}).Index(0)
	_, ok := Assert[io.Writer](v)
	assert.False(t, ok)
}
