package pretty

import (
	"encoding/hex"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/syncutil"
)

var bytesType = reflect.TypeFor[[]byte]()

// BytesHexDumpValueWriter is a [ValueWriter] that handles []byte and writes them with [hex.Dumper].
//
// It should be created with [NewBytesHexDumpValueWriter].
type BytesHexDumpValueWriter struct {
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

// NewBytesHexDumpValueWriter creates a new [BytesHexDumpValueWriter].
func NewBytesHexDumpValueWriter() *BytesHexDumpValueWriter {
	return &BytesHexDumpValueWriter{
		ShowLen:  true,
		ShowCap:  true,
		ShowAddr: false,
		MaxLen:   0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *BytesHexDumpValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if v.Type() != bytesType {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	b := v.Bytes()
	writeBytesHexDumpCommon(w, st, v, b, vw.ShowLen, vw.ShowCap, vw.ShowAddr, vw.MaxLen)
	return true
}

// Bytesable is an interface that can return a []byte.
type Bytesable interface {
	Bytes() []byte
}

var bytesableType = reflect.TypeFor[Bytesable]()

// BytesableHexDumpValueWriter is a [ValueWriter] that handles [Bytesable] and writes thems with [hex.Dumper].
//
// It should be created with [NewBytesableHexDumpValueWriter].
type BytesableHexDumpValueWriter struct {
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

// NewBytesableHexDumpValueWriter creates a new [BytesableHexDumpValueWriter].
func NewBytesableHexDumpValueWriter() *BytesableHexDumpValueWriter {
	return &BytesableHexDumpValueWriter{
		ShowLen:  true,
		ShowCap:  true,
		ShowAddr: false,
		MaxLen:   0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *BytesableHexDumpValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if !v.Type().Implements(bytesableType) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if v.Type() == reflectValueType {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	br := v.Interface().(Bytesable) //nolint:forcetypeassert // Checked above.
	b := br.Bytes()
	writeArrowWrappedString(w, ".Bytes() ")
	if b == nil {
		writeNil(w)
		return true
	}
	writeBytesHexDumpCommon(w, st, reflect.ValueOf(b), b, vw.ShowLen, vw.ShowCap, vw.ShowAddr, vw.MaxLen)
	return true
}

func writeBytesHexDumpCommon(w io.Writer, st State, v reflect.Value, b []byte, showLen bool, showCap bool, showAddr bool, maxLen int) {
	infos{
		showLen:  showLen,
		len:      len(b),
		showCap:  showCap,
		cap:      cap(b),
		showAddr: showAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.write(w)
	truncated := false
	if maxLen > 0 && len(b) > maxLen {
		b = b[:maxLen]
		truncated = true
	}
	writeString(w, "\n")
	st.IndentLevel++
	iw := GetIndentWriter(w, st.IndentString, st.IndentLevel, false)
	defer ReleaseIndentWriter(iw)
	e := getHexDumperPoolEntry(iw)
	defer releaseHexDumperPoolEntry(e)
	d := e.dumper
	mustWrite(d.Write(b))
	must(d.Close())
	if truncated {
		st.writeIndent(w)
		writeTruncated(w)
	}
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

var hexDumperPool = syncutil.PoolFor[*hexDumperPoolEntry]{
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