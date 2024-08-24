// Package pierrreerrors provides an integration with github.com/pierrre/errors.
package pierrreerrors

import (
	"io"

	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/pretty"
)

// ConfigureDefault configures [pretty.DefaultCommonValueWriter].
//
// It calls [ConfigureCommonValueWriter] with [pretty.DefaultCommonValueWriter].
func ConfigureDefault() {
	ConfigureCommonValueWriter(pretty.DefaultCommonValueWriter)
}

// Configure configures a [pretty.CommonValueWriter].
//
// It calls [ConfigureErrorValueWriter] with [pretty.CommonValueWriter.Error].
func ConfigureCommonValueWriter(vw *pretty.CommonValueWriter) {
	ConfigureErrorValueWriter(vw.Error)
}

// ConfigureErrorValueWriter configures a [pretty.ErrorValueWriter] with [Write].
func ConfigureErrorValueWriter(vw *pretty.ErrorValueWriter) {
	vw.Write = Write
}

// Write writes the error with [errverbose.Write].
func Write(w io.Writer, st pretty.State, err error) {
	st.IndentLevel++
	iw := pretty.GetIndentWriter(w, st.IndentString, st.IndentLevel, true)
	defer pretty.ReleaseIndentWriter(iw)
	errverbose.Write(iw, err)
}
