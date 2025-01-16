package pretty

import (
	"io"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/write"
)

func writeArrow(w io.Writer) {
	write.MustString(w, "=> ")
}

func writeArrowWrappedString(w io.Writer, s string) {
	writeArrow(w)
	write.MustString(w, s)
	writeArrow(w)
}

func writeTruncated(w io.Writer) {
	write.MustString(w, "<truncated>")
}

var bytesPool = syncutil.Pool[*[]byte]{
	New: func() *[]byte {
		b := make([]byte, 1024)
		return &b
	},
}
