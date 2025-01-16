package pretty

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/pretty/internal/must"
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
var DefaultPrinter = NewPrinter(DefaultConfig, DefaultCommonValueWriter)

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

// Write writes the value to the [io.Writer].
func (p *Printer) Write(w io.Writer, vi any, opts ...Option) {
	v := reflect.ValueOf(vi)
	if !v.IsValid() {
		writeNil(w)
		return
	}
	st := newState(w, p.Config.Indent)
	defer st.release()
	for _, opt := range opts {
		opt(st)
	}
	must.Handle(p.ValueWriter.WriteValue(st, v))
}

var bufPool = &bufpool.Pool{
	MaxCap: -1,
}

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

// Formatter returns a [fmt.Formatter] for the value.
func (p *Printer) Formatter(vi any, opts ...Option) fmt.Formatter {
	return &formatter{
		printer: p,
		value:   vi,
	}
}

// Option represents an option for the [Printer].
type Option func(*State)

type formatter struct {
	printer *Printer
	value   any
	opts    []Option
}

func (ft *formatter) Format(f fmt.State, verb rune) {
	ft.printer.Write(f, ft.value, ft.opts...)
}
