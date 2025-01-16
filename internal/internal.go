package internal

import (
	"io"

	"github.com/pierrre/go-libs/unsafeio"
	"github.com/pierrre/pretty/internal/must"
)

// MustWrite panics if err is not nil.
func MustWrite(_ int, err error) {
	must.NoError(err)
}

// WriteString writes a string to a [io.Writer].
func WriteString(w io.Writer, s string) (int, error) {
	return unsafeio.WriteString(w, s) //nolint:wrapcheck // The error is not wrapped.
}

// MustWriteString writes a string to a [io.Writer] and panics if an error occurs.
func MustWriteString(w io.Writer, s string) {
	MustWrite(WriteString(w, s))
}
