package pretty

import (
	"bytes"
	"io"

	"github.com/pierrre/go-libs/syncutil"
)

// IndentWriter is a [io.Writer] that indents.
//
// It is exposed for internal use.
// It should not be used outside of this package.
type IndentWriter struct {
	writer   io.Writer
	indent   string
	level    int
	indented bool
}

// NewIndentWriter creates a new [IndentWriter].
func NewIndentWriter(w io.Writer, indent string, level int, indented bool) *IndentWriter {
	iw := &IndentWriter{}
	iw.init(w, indent, level, indented)
	return iw
}

func (iw *IndentWriter) init(w io.Writer, indent string, level int, indented bool) {
	iw.writer = w
	iw.indent = indent
	iw.level = level
	iw.indented = indented
}

func (iw *IndentWriter) reset() {
	iw.writer = nil
	iw.indent = ""
	iw.level = 0
	iw.indented = false
}

// Write implements [io.Writer].
func (iw *IndentWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if !iw.indented {
			nn, err := writeIndentErr(iw.writer, iw.indent, iw.level)
			n += nn
			if err != nil {
				return n, err
			}
			iw.indented = true
		}
		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			i = len(p)
		} else {
			i++
			iw.indented = false
		}
		nn, err := iw.writer.Write(p[:i])
		n += nn
		if err != nil {
			return n, err //nolint:wrapcheck // The error is not wrapped.
		}
		p = p[i:]
	}
	return n, nil
}

var indentWriterPool = syncutil.PoolFor[*IndentWriter]{
	New: func() *IndentWriter {
		return &IndentWriter{}
	},
}

// GetIndentWriter returns a [IndentWriter] from the pool.
//
// The caller must call [IndentWriter.Release] after using it.
func GetIndentWriter(w io.Writer, indent string, level int, indented bool) *IndentWriter {
	iw := indentWriterPool.Get()
	iw.init(w, indent, level, indented)
	return iw
}

// Release releases the [IndentWriter] to the pool.
func ReleaseIndentWriter(iw *IndentWriter) {
	iw.reset()
	indentWriterPool.Put(iw)
}

const (
	indentDefault           = "\t"
	indentDefaultBytesLevel = 1000
)

var (
	indentDefaultBytes = bytes.Repeat([]byte(indentDefault), indentDefaultBytesLevel)
	indentCache        syncutil.MapFor[string, []byte]
)

func getIndentBytes(indent string, level int) []byte {
	l := len(indent) * level
	if indent == indentDefault && level <= indentDefaultBytesLevel {
		return indentDefaultBytes[:l]
	}
	b, _ := indentCache.Load(indent)
	if len(b) >= l {
		return b[:l]
	}
	b = bytes.Repeat([]byte(indent), level)
	indentCache.Store(indent, b)
	return b
}

func writeIndent(w io.Writer, indent string, level int) {
	mustWrite(writeIndentErr(w, indent, level))
}

func writeIndentErr(w io.Writer, indent string, level int) (int, error) {
	if level <= 0 {
		return 0, nil
	}
	if level == 1 {
		return writeStringErr(w, indent)
	}
	return w.Write(getIndentBytes(indent, level)) //nolint:wrapcheck // The error is not wrapped.
}
