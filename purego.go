//go:build purego

package pretty

import "io"

func writeStringErr(w io.Writer, s string) (int, error) {
	return io.WriteString(w, s)
}
