package pretty

import (
	"io"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/go-libs/unsafeio"
	"github.com/pierrre/pretty/internal"
)

func writeArrow(w io.Writer) {
	writeString(w, "=> ")
}

func writeArrowWrappedString(w io.Writer, s string) {
	writeArrow(w)
	writeString(w, s)
	writeArrow(w)
}

func writeTruncated(w io.Writer) {
	writeString(w, "<truncated>")
}

func writeString(w io.Writer, s string) {
	internal.MustWrite(writeStringErr(w, s))
}

func writeStringErr(w io.Writer, s string) (int, error) {
	return unsafeio.WriteString(w, s) //nolint:wrapcheck // The error is not wrapped.
}

func writeQuote(w io.Writer, s string) {
	internal.MustWrite(strconvio.WriteQuote(w, s))
}

var bytesPool = syncutil.Pool[*[]byte]{
	New: func() *[]byte {
		b := make([]byte, 1024)
		return &b
	},
}
