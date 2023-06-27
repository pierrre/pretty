package pierrreerrors

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/pretty"
)

func TestValueWriter(t *testing.T) {
	c := pretty.NewConfig()
	c.ValueWriters = nil
	Configure(c)
	err := errors.New("error")
	s := c.String(err)
	assert.StringHasPrefix(t, s, "(*errstack.stack) error\n\tstack\n")
}

func TestValueWriterUnexported(t *testing.T) {
	c := pretty.NewConfig()
	c.ValueWriters = nil
	Configure(c)
	err := errors.New("error")
	v := &testUnexported{v: err}
	s := c.String(v)
	assert.StringHasPrefix(t, s, "(*pierrreerrors.testUnexported) => (pierrreerrors.testUnexported) {\n\tv: (*errstack.stack) => (errstack.stack) {\n\t\terror: (*errbase.base) => (errbase.base)")
}

type testUnexported struct {
	v any
}
