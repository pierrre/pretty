package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/write"
)

var errorImplementsCache = reflectutil.NewImplementsCacheFor[error]()

// ErrorValueWriter is a [ValueWriter] that handles errors.
//
// It should be created with [NewErrorValueWriter].
type ErrorValueWriter struct {
	// Write writes the error.
	// Default: [ErrorValueWriter.WriteError].
	Write func(st *State, err error)
}

// NewErrorValueWriter creates a new [ErrorValueWriter] with default values.
func NewErrorValueWriter() *ErrorValueWriter {
	vw := &ErrorValueWriter{}
	vw.Write = vw.WriteError
	return vw
}

// WriteValue implements [ValueWriter].
func (vw *ErrorValueWriter) WriteValue(st *State, v reflect.Value) bool {
	typ := v.Type()
	if !errorImplementsCache.ImplementedBy(typ) {
		return false
	}
	err, ok := itfassert.Assert[error](v)
	if !ok {
		return false
	}
	writeArrowWrappedString(st.Writer, ".Error() ")
	vw.Write(st, err)
	return true
}

// WriteError writes the error with error.Error.
func (vw *ErrorValueWriter) WriteError(st *State, err error) {
	write.Must(strconvio.WriteQuote(st.Writer, err.Error()))
}
