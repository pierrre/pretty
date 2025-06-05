package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

// StringValueWriter is a [ValueWriter] that handles string values.
//
// It should be created with [NewStringValueWriter].
type StringValueWriter struct {
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

// NewStringValueWriter creates a new [StringValueWriter] with default values.
func NewStringValueWriter() *StringValueWriter {
	return &StringValueWriter{
		ShowLen:  true,
		ShowAddr: false,
		Quote:    true,
		MaxLen:   0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StringValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !vw.supportsType(v.Type()) {
		return false
	}
	return vw.writeValue(st, v)
}

func (vw *StringValueWriter) writeValue(st *State, v reflect.Value) bool {
	s := v.String()
	writeStringValue(st, s, vw.ShowLen, vw.ShowAddr, uintptr(v.UnsafePointer()), vw.Quote, vw.MaxLen)
	return true
}

// SupportsType implements [TypeSupportChecker].
func (vw *StringValueWriter) SupportsType(typ reflect.Type) ValueWriterFunc {
	if vw.supportsType(typ) {
		return vw.writeValue
	}
	return nil
}

func (vw *StringValueWriter) supportsType(typ reflect.Type) bool {
	return typ.Kind() == reflect.String
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
