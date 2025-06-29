package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

var errorImplementsCache = reflectutil.NewImplementsCacheFor[error]()

// ErrorWriter is a [ValueWriter] that handles errors.
//
// It should be created with [NewErrorWriter].
type ErrorWriter struct {
	ValueWriter
}

// NewErrorWriter creates a new [ErrorWriter] with default values.
func NewErrorWriter(vw ValueWriter) *ErrorWriter {
	return &ErrorWriter{
		ValueWriter: vw,
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
	write.MustString(st.Writer, "Error: ")
	write.Must(strconvio.WriteQuote(st.Writer, err.Error()))
	write.MustString(st.Writer, ",\n")
	switch err := err.(type) { //nolint:errorlint // We want to check which interface is implemented by the current error.
	case interface{ Unwrap() error }:
		st.WriteIndent()
		write.MustString(st.Writer, "Unwrap: ")
		st.KnownType = false // We want to show the type of the unwrapped error.
		must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(err.Unwrap())))
		write.MustString(st.Writer, ",\n")
	case interface{ Unwrap() []error }:
		st.WriteIndent()
		write.MustString(st.Writer, "Unwrap: ")
		st.KnownType = false // We want to show the type of the unwrapped errors.
		must.Handle(vw.ValueWriter.WriteValue(st, reflect.ValueOf(err.Unwrap())))
		write.MustString(st.Writer, ",\n")
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
