package write_test

import (
	"errors"
	"io"
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

func TestMustString(t *testing.T) {
	MustString(io.Discard, "test")
}
