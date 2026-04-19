package pretty

import (
	"reflect"
)

// Vars returns a slice of variadic arguments.
// It allows calling a [Printer] with the result of a function returning multiple values.
func Vars(vs ...any) []any {
	return vs
}

func writeArrow(st *State) {
	st.Writer.AppendString("=> ")
}

func writeArrowWrappedString(st *State, s string) {
	writeArrow(st)
	st.Writer.AppendString(s)
	writeArrow(st)
}

func writeTruncated(st *State) {
	st.Writer.AppendString("<truncated>")
}

func callSupportCheckerPointer[P interface {
	*T
	SupportChecker
}, T any](p P, typ reflect.Type) ValueWriter {
	if p != nil {
		return p.Supports(typ)
	}
	return nil
}
