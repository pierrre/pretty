//go:build purego

package pretty

import "io"

func writeString(w io.Writer, s string) (int, error) {
	return io.WriteString(w, s)
}
