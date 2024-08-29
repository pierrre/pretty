package indent

import (
	"bytes"
	"io"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal"
)

const (
	// Default is the default string used for indentation.
	Default           = "\t"
	defaultBytesLevel = 1000
)

var (
	defaultBytes = bytes.Repeat([]byte(Default), defaultBytesLevel)
	bytesCache   syncutil.MapFor[string, []byte]
)

func getBytes(str string, level int) []byte {
	l := len(str) * level
	if str == Default && level <= defaultBytesLevel {
		return defaultBytes[:l]
	}
	b, _ := bytesCache.Load(str)
	if len(b) >= l {
		return b[:l]
	}
	b = bytes.Repeat([]byte(str), level)
	bytesCache.Store(str, b)
	return b
}

// Write writes an indentation string to a [io.Writer].
func Write(w io.Writer, str string, level int) (int, error) {
	if level <= 0 {
		return 0, nil
	}
	if level == 1 {
		return internal.WriteString(w, str) //nolint:wrapcheck // The error is not wrapped.
	}
	return w.Write(getBytes(str, level)) //nolint:wrapcheck // The error is not wrapped.
}

// MustWrite writes an indentation string to a [io.Writer] and panics if an error occurs.
func MustWrite(w io.Writer, str string, level int) {
	internal.MustWrite(Write(w, str, level))
}

// Writer is a [io.Writer] that indents.
type Writer struct {
	writer   io.Writer
	string   string
	level    int
	indented bool
}

var writerPool = syncutil.PoolFor[*Writer]{
	New: func() *Writer {
		return &Writer{}
	},
}

// NewWriter creates a new [Writer].
func NewWriter(w io.Writer, str string, level int, indented bool) *Writer {
	iw := writerPool.Get()
	iw.writer = w
	iw.string = str
	iw.level = level
	iw.indented = indented
	return iw
}

// Write implements [io.Writer].
func (iw *Writer) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if !iw.indented {
			nn, err := Write(iw.writer, iw.string, iw.level)
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

// Release releases the [Writer].
//
// It must not be used after calling this method.
func (iw *Writer) Release() {
	iw.writer = nil
	iw.string = ""
	iw.level = 0
	iw.indented = false
	writerPool.Put(iw)
}
