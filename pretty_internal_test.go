package pretty

import (
	"io"
)

func WriteIndent(w io.Writer, indent string, level int) {
	writeIndent(w, indent, level)
}
