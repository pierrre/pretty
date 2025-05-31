// Package write provides helpers to write.
//
//nolint:gosec // It uses unsafe.
package write

import (
	"io"
	"unsafe" //nolint:depguard // The current package is unsafe.
)

// Must panics if err is not nil.
func Must(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

// String writes a string to a [io.Writer].
func String(w io.Writer, s string) (int, error) {
	return w.Write( //nolint:gosec,wrapcheck // The error is not wrapped.
		unsafe.Slice(
			unsafe.StringData(s),
			len(s),
		),
	)
}

// MustString writes a string to a [io.Writer] and panics if an error occurs.
func MustString(w io.Writer, s string) {
	_, err := w.Write(
		unsafe.Slice(
			unsafe.StringData(s),
			len(s),
		),
	)
	if err != nil {
		panic(err)
	}
}
