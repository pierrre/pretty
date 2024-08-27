package pretty

import (
	"io"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/go-libs/unsafeio"
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
	mustWrite(writeStringErr(w, s))
}

func writeStringErr(w io.Writer, s string) (int, error) {
	return unsafeio.WriteString(w, s) //nolint:wrapcheck // The error is not wrapped.
}

func writeQuote(w io.Writer, s string) {
	mustWrite(strconvio.WriteQuote(w, s))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustWrite(_ int, err error) {
	must(err)
}
