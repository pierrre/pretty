package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// StringWriter is a [ValueWriter] that handles string values.
//
// It should be created with [NewStringWriter].
type StringWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowAddr shows the address.
	// Default: false.
	ShowAddr bool
	// Quote quotes the string.
	// Default: true.
	Quote bool
	// MaxLen is the maximum length of the string.
	// Default: 0 (no limit).
	MaxLen int
}

// NewStringWriter creates a new [StringWriter] with default values.
func NewStringWriter() *StringWriter {
	return &StringWriter{
		ShowLen:  true,
		ShowAddr: false,
		Quote:    true,
		MaxLen:   0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StringWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.String {
		return false
	}
	writeStringValue(st, v.String(), vw.ShowLen, vw.ShowAddr, uintptr(v.UnsafePointer()), vw.Quote, vw.MaxLen)
	return true
}

// Supports implements [SupportChecker].
func (vw *StringWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.String {
		res = vw
	}
	return res
}

func writeStringValue(st *State, s string, showLen bool, showAddr bool, addr uintptr, quote bool, maxLen int) {
	infos{
		showLen:  showLen,
		len:      len(s),
		showAddr: showAddr,
		addr:     addr,
	}.writeWithTrailingSpace(st)
	truncated := false
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
		truncated = true
	}
	if quote {
		write.Must(strconvio.WriteQuote(st.Writer, s))
	} else {
		write.MustString(st.Writer, s)
	}
	if truncated {
		write.MustString(st.Writer, " ")
		writeTruncated(st.Writer)
	}
}
