// Package indent provides utilities for writing indented output.
package indent

import (
	"bytes"
	"io"

	"github.com/pierrre/go-libs/syncutil"
)

const (
	// Default is the default string used for indentation.
	Default           = "\t"
	defaultBytesLevel = 1000
)

var (
	defaultBytes = bytes.Repeat([]byte(Default), defaultBytesLevel)
	bytesCache   syncutil.Map[string, []byte]
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

// Append appends an indentation string to a []byte and returns the result.
func Append(dst []byte, str string, level int) []byte {
	if level <= 0 {
		return dst
	}
	if level == 1 {
		return append(dst, str...)
	}
	return append(dst, getBytes(str, level)...)
}

// Writer is a [io.Writer] that indents.
type Writer struct {
	writer   io.Writer
	bytes    []byte
	indented bool
}

var writerPool = syncutil.Pool[*Writer]{
	New: func() *Writer {
		return &Writer{}
	},
}

// NewWriter creates a new [Writer].
func NewWriter(w io.Writer, str string, level int, indented bool) *Writer {
	iw := writerPool.Get()
	iw.writer = w
	iw.bytes = getBytes(str, level)
	iw.indented = indented
	return iw
}

// Write implements [io.Writer].
func (iw *Writer) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if !iw.indented {
			if len(iw.bytes) != 0 {
				nn, err := iw.writer.Write(iw.bytes)
				n += nn
				if err != nil {
					return n, err //nolint:wrapcheck // The error is not wrapped.
				}
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
	iw.bytes = nil
	iw.indented = false
	writerPool.Put(iw)
}
