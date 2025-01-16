package internal

import (
	"io"

	"github.com/pierrre/go-libs/unsafeio"
)

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustWrite panics if err is not nil.
func MustWrite(_ int, err error) {
	Must(err)
}

// WriteString writes a string to a [io.Writer].
func WriteString(w io.Writer, s string) (int, error) {
	return unsafeio.WriteString(w, s) //nolint:wrapcheck // The error is not wrapped.
}

// MustWriteString writes a string to a [io.Writer] and panics if an error occurs.
func MustWriteString(w io.Writer, s string) {
	MustWrite(WriteString(w, s))
}

// MustHandle panics if h is false.
func MustHandle(h bool) {
	if !h {
		panic("not handled")
	}
}
