package pretty

import (
	"io"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal"
)

func writeArrow(w io.Writer) {
	internal.MustWriteString(w, "=> ")
}

func writeArrowWrappedString(w io.Writer, s string) {
	writeArrow(w)
	internal.MustWriteString(w, s)
	writeArrow(w)
}

func writeTruncated(w io.Writer) {
	internal.MustWriteString(w, "<truncated>")
}

var bytesPool = syncutil.Pool[*[]byte]{
	New: func() *[]byte {
		b := make([]byte, 1024)
		return &b
	},
}
