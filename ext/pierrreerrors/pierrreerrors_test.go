package pierrreerrors

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/pretty"
)

func init() {
	ConfigureDefault()
	pretty.DefaultCommonWriter.CanInterface = nil
}

func TestValueWriter(t *testing.T) {
	err := errors.New("error")
	s := pretty.String(err)
	assert.StringHasPrefix(t, s, "[*github.com/pierrre/errors/errstack.stack] => .Error() => error\n\tstack:\n")
}

func TestValueWriterNil(t *testing.T) {
	var err error = (*testError)(nil)
	s := pretty.String(err)
	assert.Equal(t, s, "[*github.com/pierrre/pretty/ext/pierrreerrors.testError] <nil>")
}

func TestValueWriterUnexported(t *testing.T) {
	err := errors.New("error")
	v := &testUnexported{v: err}
	s := pretty.String(v)
	assert.StringHasPrefix(t, s, "[*github.com/pierrre/pretty/ext/pierrreerrors.testUnexported] => {\n\tv: [*github.com/pierrre/errors/errstack.stack] => {\n\t\terror: [*errors.errorString] => {")
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}

type testUnexported struct {
	v any
}
