package pretty

import (
	"io"
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/indent"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

var errorImplementsCache = reflectutil.NewImplementsCacheFor[error]()

// ErrorWriter is a [ValueWriter] that handles errors.
//
// It writes the error's message, then the custom function, then the unwrapped error(s) if any.
//
// It should be created with [NewErrorWriter].
type ErrorWriter struct {
	ValueWriter
	// Writers is a list of custom functions that are called when an error is written.
	// Default: {[WriteVerboseError]}.
	Writers []func(*State, error)
}

// NewErrorWriter creates a new [ErrorWriter] with default values.
func NewErrorWriter(vw ValueWriter) *ErrorWriter {
	return &ErrorWriter{
		ValueWriter: vw,
		Writers: []func(*State, error){
			WriteVerboseError,
		},
	}
}

// WriteValue implements [ValueWriter].
func (vw *ErrorWriter) WriteValue(st *State, v reflect.Value) bool {
	if !errorImplementsCache.ImplementedBy(v.Type()) {
		return false
	}
	err, ok := itfassert.Assert[error](v)
	if !ok {
		return false
	}
	write.MustString(st.Writer, "{\n")
	st.IndentLevel++
	st.WriteIndent()
	write.MustString(st.Writer, "Error(): ")
	write.Must(strconvio.WriteQuote(st.Writer, err.Error()))
	write.MustString(st.Writer, ",\n")
	for _, w := range vw.Writers {
		w(st, err)
	}
	switch err := err.(type) { //nolint:errorlint // We want to check which interface is implemented by the current error.
	case interface{ Unwrap() error }:
		e := err.Unwrap()
		if e != nil {
			st.WriteIndent()
			write.MustString(st.Writer, "Unwrap(): ")
			st.KnownType = false // We want to show the type of the unwrapped error.
			must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(e)))
			write.MustString(st.Writer, ",\n")
		}
	case interface{ Unwrap() []error }:
		errs := err.Unwrap()
		if len(errs) > 0 {
			st.WriteIndent()
			write.MustString(st.Writer, "Unwrap(): ")
			st.KnownType = false // We want to show the type of the unwrapped errors.
			must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(errs)))
			write.MustString(st.Writer, ",\n")
		}
	}
	st.IndentLevel--
	st.WriteIndent()
	write.MustString(st.Writer, "}")
	return true
}

// Supports implements [SupportChecker].
func (vw *ErrorWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if errorImplementsCache.ImplementedBy(typ) {
		res = vw
	}
	return res
}

// VerboseError is an interface that can be implemented by errors to provide a verbose error message.
type VerboseError interface {
	// ErrorVerbose writes the error verbose message.
	// It must only write the verbose message of the error, not the error chain.
	ErrorVerbose(w io.Writer)
}

// WriteVerboseError writes the verbose error message of an error that implements [VerboseError].
func WriteVerboseError(st *State, err error) {
	v, ok := err.(VerboseError)
	if !ok {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "ErrorVerbose(): ")
	st.IndentLevel++
	iw := indent.NewWriter(st.Writer, st.IndentString, st.IndentLevel, true)
	v.ErrorVerbose(iw)
	write.MustString(iw, ",\n")
	iw.Release()
	st.IndentLevel--
}
