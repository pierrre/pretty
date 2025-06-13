package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/write"
)

var errorImplementsCache = reflectutil.NewImplementsCacheFor[error]()

// ErrorWriter is a [ValueWriter] that handles errors.
//
// It should be created with [NewErrorWriter].
type ErrorWriter struct {
	// Write writes the error.
	// Default: [ErrorWriter.WriteError].
	Write func(st *State, err error)
}

// NewErrorWriter creates a new [ErrorWriter] with default values.
func NewErrorWriter() *ErrorWriter {
	vw := &ErrorWriter{}
	vw.Write = vw.WriteError
	return vw
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
	writeArrowWrappedString(st.Writer, ".Error() ")
	vw.Write(st, err)
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

// WriteError writes the error with error.Error.
func (vw *ErrorWriter) WriteError(st *State, err error) {
	write.Must(strconvio.WriteQuote(st.Writer, err.Error()))
}
