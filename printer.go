package pretty

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/bufpool"
)

// Write writes the value to the [io.Writer] with [DefaultPrinter].
func Write(w io.Writer, vi any, opts ...Option) {
	DefaultPrinter.Write(w, vi, opts...)
}

// String returns the value as a string with [DefaultPrinter].
func String(vi any, opts ...Option) string {
	return DefaultPrinter.String(vi, opts...)
}

// Formatter returns a [fmt.Formatter] for the value with [DefaultPrinter].
func Formatter(vi any, opts ...Option) fmt.Formatter {
	return DefaultPrinter.Formatter(vi, opts...)
}

// DefaultPrinter is the default [Printer].
var DefaultPrinter = NewPrinterCommon(DefaultConfig, DefaultCommonValueWriter)

// Printer pretty-prints values.
//
// It should be created with [NewPrinter].
type Printer struct {
	Config      *Config
	ValueWriter ValueWriter
}

// NewPrinter creates a new [Printer].
func NewPrinter(c *Config, vw ValueWriter) *Printer {
	return &Printer{
		Config:      c,
		ValueWriter: vw,
	}
}

// NewPrinterCommon creates a new [Printer] with a [CommonValueWriter].
//
// It calls [CommonValueWriter.ConfigureWithPrinter] with the created [Printer].
func NewPrinterCommon(c *Config, vw *CommonValueWriter) *Printer {
	p := NewPrinter(c, vw.WriteValue)
	vw.ConfigureWithPrinter(p)
	return p
}

// Write writes the value to the [io.Writer].
func (p *Printer) Write(w io.Writer, vi any, opts ...Option) {
	v := reflect.ValueOf(vi)
	if !v.IsValid() {
		writeNil(w)
		return
	}
	st := getState()
	defer st.release()
	st.IndentString = p.Config.Indent
	for _, opt := range opts {
		st = opt(st)
	}
	mustHandle(p.ValueWriter(w, st, v))
}

var bufPool = &bufpool.Pool{}

// String returns the value as a string.
func (p *Printer) String(vi any, opts ...Option) string {
	buf := p.getBuf(vi, opts...)
	defer bufPool.Put(buf)
	return buf.String()
}

func (p *Printer) getBuf(vi any, opts ...Option) *bytes.Buffer {
	buf := bufPool.Get()
	p.Write(buf, vi, opts...)
	return buf
}

func (p *Printer) compare(a, b reflect.Value) int {
	aBuf := p.getBuf(a)
	defer bufPool.Put(aBuf)
	bBuf := p.getBuf(b)
	defer bufPool.Put(bBuf)
	return bytes.Compare(aBuf.Bytes(), bBuf.Bytes())
}

// Formatter returns a [fmt.Formatter] for the value.
func (p *Printer) Formatter(vi any, opts ...Option) fmt.Formatter {
	return &formatter{
		printer: p,
		value:   vi,
	}
}

// Option represents an option for the [Printer].
type Option func(State) State

type formatter struct {
	printer *Printer
	value   any
	opts    []Option
}

func (ft *formatter) Format(f fmt.State, verb rune) {
	ft.printer.Write(f, ft.value, ft.opts...)
}
