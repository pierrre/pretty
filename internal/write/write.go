// Package write provides helpers to write.
package write

import (
	"io"

	"github.com/pierrre/go-libs/unsafeio"
	"github.com/pierrre/pretty/internal/must"
)

// Must panics if err is not nil.
func Must(_ int, err error) {
	must.NoError(err)
}

// String writes a string to a [io.Writer].
func String(w io.Writer, s string) (int, error) {
	return unsafeio.WriteString(w, s) //nolint:wrapcheck // The error is not wrapped.
}

// MustString writes a string to a [io.Writer] and panics if an error occurs.
func MustString(w io.Writer, s string) {
	Must(
		String(w, s),
	)
}
