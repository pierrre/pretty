package pretty

import (
	"reflect"
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
	if v.Kind() != reflect.String {
		return false
	}
	s := v.String()
	writeStringValue(st, s, vw.ShowLen, vw.ShowAddr, uintptr(v.UnsafePointer()), vw.Quote, vw.MaxLen)
	return true
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
		writeQuote(st.Writer, s)
	} else {
		writeString(st.Writer, s)
	}
	if truncated {
		writeString(st.Writer, " ")
		writeTruncated(st.Writer)
	}
}
