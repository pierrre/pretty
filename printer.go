package pretty

import (
	"fmt"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// Write writes the value to the [io.Writer] with [DefaultPrinter].
func Write(w io.Writer, vi any) {
	DefaultPrinter.Write(w, vi)
}

// String returns the value as a string with [DefaultPrinter].
func String(vi any) string {
	return DefaultPrinter.String(vi)
}

// Formatter returns a [fmt.Formatter] for the value with [DefaultPrinter].
func Formatter(vi any) fmt.Formatter {
	return DefaultPrinter.Formatter(vi)
}

// DefaultPrinter is the default [Printer].
var DefaultPrinter = NewPrinter(DefaultCommonValueWriter)

// Printer pretty-prints values.
//
// It should be created with [NewPrinter].
type Printer struct {
	// ValueWriter is the [ValueWriter] used to write values.
	ValueWriter ValueWriter
	// PanicRecover indicates whether to recover from panics.
	// Default: true.
	PanicRecover bool
	// Indent is the string used to indent.
	// Default: "\t".
	Indent string
}

// NewPrinter creates a new [Printer].
func NewPrinter(vw ValueWriter) *Printer {
	return &Printer{
		ValueWriter:  vw,
		PanicRecover: true,
		Indent:       "\t",
	}
}

// Write writes the value to the [io.Writer].
func (p *Printer) Write(w io.Writer, vi any) {
	if p.PanicRecover {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			writePanic(w, r)
		}()
	}
	v := reflect.ValueOf(vi)
	if checkInvalidNil(w, v) {
		return
	}
	st := newState(w, p.Indent)
	defer st.release()
	must.Handle(p.ValueWriter.WriteValue(st, v))
}

func writePanic(w io.Writer, r any) {
	_, _ = write.String(w, "<panic>: ")
	_, _ = fmt.Fprint(w, r)
}

var bufPool = &bufpool.Pool{
	MaxCap: -1,
}

// String returns the value as a string.
func (p *Printer) String(vi any) string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	p.Write(buf, vi)
	return buf.String()
}

// Formatter returns a [fmt.Formatter] for the value.
func (p *Printer) Formatter(vi any) fmt.Formatter {
	return &formatter{
		printer: p,
		value:   vi,
	}
}

type formatter struct {
	printer *Printer
	value   any
}

func (ft *formatter) Format(f fmt.State, verb rune) {
	ft.printer.Write(f, ft.value)
}
