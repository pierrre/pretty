package internal_test

import (
	"errors"
	"io"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty/internal"
)

func TestMustWriteNoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		MustWrite(0, nil)
	})
}

func TestMustWritePanic(t *testing.T) {
	assert.Panics(t, func() {
		MustWrite(0, errors.New("test"))
	})
}

func TestMustWriteString(t *testing.T) {
	MustWriteString(io.Discard, "test")
}
