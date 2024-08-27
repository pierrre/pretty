package pretty

import (
	"io"
	"reflect"
)

var errorType = reflect.TypeFor[error]()

// ErrorValueWriter is a [ValueWriter] that handles errors.
//
// It should be created with [NewErrorValueWriter].
type ErrorValueWriter struct {
	// Write writes the error.
	// Default: [ErrorValueWriter.WriteError].
	Write func(w io.Writer, st State, err error)
}

// NewErrorValueWriter creates a new [ErrorValueWriter] with default values.
func NewErrorValueWriter() *ErrorValueWriter {
	vw := &ErrorValueWriter{}
	vw.Write = vw.WriteError
	return vw
}

// WriteValue implements [ValueWriter].
func (vw *ErrorValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
	if !v.Type().Implements(errorType) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	writeArrowWrappedString(w, ".Error() ")
	vw.Write(w, st, err)
	return true
}

// WriteError writes the error with error.Error.
func (vw *ErrorValueWriter) WriteError(w io.Writer, st State, err error) {
	writeQuote(w, err.Error())
}
