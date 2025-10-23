package pretty

import (
	"encoding/hex"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/indent"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

var bytesType = reflect.TypeFor[[]byte]()

// BytesHexDumpWriter is a [ValueWriter] that handles []byte and writes them with [hex.Dumper].
//
// It should be created with [NewBytesHexDumpWriter].
type BytesHexDumpWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// ShowAddr shows the address.
	// Default: false.
	ShowAddr bool
	// MaxLen is the maximum length of the bytes.
	// Default: 0 (no limit).
	MaxLen int
}

// NewBytesHexDumpWriter creates a new [BytesHexDumpWriter].
func NewBytesHexDumpWriter() *BytesHexDumpWriter {
	return &BytesHexDumpWriter{
		ShowLen:  true,
		ShowCap:  true,
		ShowAddr: false,
		MaxLen:   0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *BytesHexDumpWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Slice || v.Type() != bytesType {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	b := v.Bytes()
	writeBytesHexDumpCommon(st, v, b, vw.ShowLen, vw.ShowCap, vw.ShowAddr, vw.MaxLen)
	return true
}

// Supports implements [SupportChecker].
func (vw *BytesHexDumpWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Slice && typ == bytesType {
		res = vw
	}
	return res
}

// Bytesable is an interface that can return a []byte.
type Bytesable interface {
	Bytes() []byte
}

var bytesableImplementsCache = reflectutil.NewImplementsCacheFor[Bytesable]()

// BytesableHexDumpWriter is a [ValueWriter] that handles [Bytesable] and writes thems with [hex.Dumper].
//
// If [Bytesable.Bytes] panics, [BytesableHexDumpWriter.WriteValue] returns false.
//
// It should be created with [NewBytesableHexDumpWriter].
type BytesableHexDumpWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// ShowAddr shows the address.
	// Default: false.
	ShowAddr bool
	// MaxLen is the maximum length of the bytes.
	// Default: 0 (no limit).
	MaxLen int
}

// NewBytesableHexDumpWriter creates a new [BytesableHexDumpWriter].
func NewBytesableHexDumpWriter() *BytesableHexDumpWriter {
	return &BytesableHexDumpWriter{
		ShowLen:  true,
		ShowCap:  true,
		ShowAddr: false,
		MaxLen:   0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *BytesableHexDumpWriter) WriteValue(st *State, v reflect.Value) bool {
	typ := v.Type()
	if typ == reflectValueType {
		return false
	}
	if !bytesableImplementsCache.ImplementedBy(typ) {
		return false
	}
	br, ok := itfassert.Assert[Bytesable](v)
	if !ok {
		return false
	}
	b, ok := func() (_ []byte, ok bool) {
		defer func() {
			if !ok {
				_ = recover()
			}
		}()
		return br.Bytes(), true
	}()
	if !ok {
		return false
	}
	writeArrowWrappedString(st.Writer, "Bytes() ")
	if b == nil {
		writeNil(st.Writer)
		return true
	}
	writeBytesHexDumpCommon(st, reflect.ValueOf(b), b, vw.ShowLen, vw.ShowCap, vw.ShowAddr, vw.MaxLen)
	return true
}

// Supports implements [SupportChecker].
func (vw *BytesableHexDumpWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ != reflectValueType && bytesableImplementsCache.ImplementedBy(typ) {
		res = vw
	}
	return res
}

func writeBytesHexDumpCommon(st *State, v reflect.Value, b []byte, showLen bool, showCap bool, showAddr bool, maxLen int) {
	infos{
		showLen:  showLen,
		len:      len(b),
		showCap:  showCap,
		cap:      cap(b),
		showAddr: showAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.write(st)
	truncated := false
	if maxLen > 0 && len(b) > maxLen {
		b = b[:maxLen]
		truncated = true
	}
	write.MustString(st.Writer, "\n")
	st.IndentLevel++
	iw := indent.NewWriter(st.Writer, st.IndentString, st.IndentLevel, false)
	e := getHexDumperPoolEntry(iw)
	d := e.dumper
	write.Must(d.Write(b))
	must.NoError(d.Close())
	releaseHexDumperPoolEntry(e)
	iw.Release()
	if truncated {
		st.WriteIndent()
		writeTruncated(st.Writer)
		write.MustString(st.Writer, "\n")
	}
	st.IndentLevel--
	st.WriteIndent()
}

type hexDumperPoolEntry struct {
	dumper        io.WriteCloser
	original      io.WriteCloser
	writerWrapper *writerWrapper
}

func newHexDumperPoolEntry() *hexDumperPoolEntry {
	ww := &writerWrapper{}
	return &hexDumperPoolEntry{
		dumper:        hex.Dumper(ww),
		original:      hex.Dumper(ww),
		writerWrapper: ww,
	}
}

var hexDumperPool = syncutil.Pool[*hexDumperPoolEntry]{
	New: newHexDumperPoolEntry,
}

func getHexDumperPoolEntry(w io.Writer) *hexDumperPoolEntry {
	e := hexDumperPool.Get()
	e.writerWrapper.Writer = w
	return e
}

func releaseHexDumperPoolEntry(e *hexDumperPoolEntry) {
	v1 := reflect.ValueOf(e.dumper).Elem()
	v2 := reflect.ValueOf(e.original).Elem()
	v1.Set(v2)
	e.writerWrapper.Writer = nil
	hexDumperPool.Put(e)
}

type writerWrapper struct {
	io.Writer
}
