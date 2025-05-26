// Package pierrreerrors provides an integration with github.com/pierrre/errors.
package pierrreerrors

import (
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/indent"
)

// ConfigureDefault configures [pretty.DefaultCommonValueWriter].
//
// It calls [ConfigureCommonValueWriter] with [pretty.DefaultCommonValueWriter].
func ConfigureDefault() {
	ConfigureCommonValueWriter(pretty.DefaultCommonValueWriter)
}

// ConfigureCommonValueWriter configures a [pretty.CommonValueWriter].
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
func Write(st *pretty.State, err error) {
	st.IndentLevel++
	iw := indent.NewWriter(st.Writer, st.IndentString, st.IndentLevel, true)
	errverbose.Write(iw, err)
	iw.Release()
	st.IndentLevel--
}
