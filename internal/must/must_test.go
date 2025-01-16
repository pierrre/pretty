package must_test

import (
	"errors"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty/internal/must"
)

func TestNoError(t *testing.T) {
	assert.NotPanics(t, func() {
		NoError(nil)
	})
	assert.Panics(t, func() {
		NoError(errors.New("test"))
	})
}

func TestHandle(t *testing.T) {
	assert.NotPanics(t, func() {
		Handle(true)
	})
	assert.Panics(t, func() {
		Handle(false)
	})
}
