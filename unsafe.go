//go:build !purego

package pretty

import (
	"io"
	"unsafe" //nolint:depguard // Required for string to []byte conversion.
)

func writeString(w io.Writer, s string) (int, error) {
	return w.Write(unsafeStringToBytes(s)) //nolint:wrapcheck // The error is not wrapped.
}

func unsafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
