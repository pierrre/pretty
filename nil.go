package pretty

import (
	"reflect"
)

func checkNil(st *State, v reflect.Value) bool {
	if v.IsNil() {
		writeNil(st)
		return true
	}
	return false
}

func writeNil(st *State) {
	st.Writer.AppendString("<nil>")
}
