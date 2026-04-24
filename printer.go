package pretty

import (
	"fmt"
	"io"
	"reflect"

	"github.com/pierrre/go-libs/panicutil"
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
	err := p.writeTo(w, vi)
	if err != nil {
		panic(err)
	}
}

// WriteErr writes the value to the [io.Writer] and returns an error if it occurs.
//
// It recovers from panics and returns them as errors.
func (p *Printer) WriteErr(w io.Writer, vi any) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = panicutil.NewError(r)
		}
	}()
	return p.writeTo(w, vi)
}

func (p *Printer) writeTo(w io.Writer, vi any) error {
	st := newState(p.Indent)
	defer st.release()
	p.write(st, vi)
	n, err := w.Write(st.Writer)
	if err != nil {
		return err //nolint:wrapcheck // No need to wrap error.
	}
	if n != len(st.Writer) {
		return io.ErrShortWrite
	}
	return nil
}

// String returns the value as a string.
func (p *Printer) String(vi any) string {
	st := newState(p.Indent)
	defer st.release()
	p.write(st, vi)
	return st.Writer.String()
}

func (p *Printer) write(st *State, vi any) {
	v := reflect.ValueOf(vi)
	if checkInvalidNil(st, v) {
		return
	}
	p.ValueWriter.WriteValue(st, v)
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
