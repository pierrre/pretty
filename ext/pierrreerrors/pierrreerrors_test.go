package pierrreerrors

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/pretty"
)

func init() {
	ConfigureDefault()
	pretty.DefaultCommonValueWriter.ConfigureTest()
	pretty.DefaultCommonValueWriter.CanInterface = nil
}

func TestValueWriter(t *testing.T) {
	err := errors.New("error")
	s := pretty.String(err)
	assert.StringHasPrefix(t, s, "[*errstack.stack] error\n\tstack\n")
}

func TestValueWriterNil(t *testing.T) {
	var err error = (*testError)(nil)
	s := pretty.String(err)
	assert.Equal(t, s, "[*pierrreerrors.testError] <nil>")
}

func TestValueWriterUnexported(t *testing.T) {
	err := errors.New("error")
	v := &testUnexported{v: err}
	s := pretty.String(v)
	assert.StringHasPrefix(t, s, "[*pierrreerrors.testUnexported] => {\n\tv: [*errstack.stack] => {\n\t\terror: [*errors.errorString] => {")
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}

type testUnexported struct {
	v any
}
