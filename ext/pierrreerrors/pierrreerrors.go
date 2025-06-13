// Package pierrreerrors provides an integration with github.com/pierrre/errors.
package pierrreerrors

import (
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/indent"
)

// ConfigureDefault configures [pretty.DefaultCommonWriter].
//
// It calls [ConfigureCommonWriter] with [pretty.DefaultCommonWriter].
func ConfigureDefault() {
	ConfigureCommonWriter(pretty.DefaultCommonWriter)
}

// ConfigureCommonWriter configures a [pretty.CommonWriter].
//
// It calls [ConfigureErrorWriter] with [pretty.CommonWriter.Error].
func ConfigureCommonWriter(vw *pretty.CommonWriter) {
	ConfigureErrorWriter(vw.Error)
}

// ConfigureErrorWriter configures a [pretty.ErrorWriter] with [Write].
func ConfigureErrorWriter(vw *pretty.ErrorWriter) {
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
