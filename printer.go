package pretty

import (
	"fmt"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/panicutil"
	"github.com/pierrre/pretty/internal/must"
)

// Write writes the value to the [io.Writer] with [DefaultPrinter].
func Write(w io.Writer, vi any) {
	DefaultPrinter.Write(w, vi)
}

// WriteErr writes the value to the [io.Writer] with [DefaultPrinter], and returns an error if it occurs.
func WriteErr(w io.Writer, vi any) error {
	return DefaultPrinter.WriteErr(w, vi)
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
//
// It uses [DefaultWriter].
var DefaultPrinter = NewPrinter(DefaultWriter)

// Printer pretty-prints values.
//
// It should be created with [NewPrinter].
type Printer struct {
	// ValueWriter is the [ValueWriter] used to write values.
	ValueWriter ValueWriter
	// Indent is the string used to indent.
	// Default: "\t".
	Indent string
}

// NewPrinter creates a new [Printer].
func NewPrinter(vw ValueWriter) *Printer {
	return &Printer{
		ValueWriter: vw,
		Indent:      "\t",
	}
}

// Write writes the value to the [io.Writer].
//
// It panics if there is a write error.
// For error handling, see [Printer.WriteErr].
func (p *Printer) Write(w io.Writer, vi any) {
	v := reflect.ValueOf(vi)
	if checkInvalidNil(w, v) {
		return
	}
	st := newState(w, p.Indent)
	defer st.release()
	must.Handle(p.ValueWriter.WriteValue(st, v))
}

// WriteErr writes the value to the [io.Writer] and returns an error if it occurs.
func (p *Printer) WriteErr(w io.Writer, vi any) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = panicutil.NewError(r)
		}
	}()
	p.Write(w, vi)
	return nil
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
	_ = ft.printer.WriteErr(f, ft.value)
}
