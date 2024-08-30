package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCases([]*testCase{
		{
			name:  "Invalid",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{func(st *State, v reflect.Value) bool {
					return vw.Kind.WriteValue(st, reflect.ValueOf(nil))
				}}
			},
		},
	})
}
