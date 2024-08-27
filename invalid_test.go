package pretty_test

import (
	"io"
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCases([]*testCase{
		{
			name:  "Invalid",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{func(w io.Writer, st State, v reflect.Value) bool {
					return vw.Kind.WriteValue(w, st, reflect.ValueOf(nil))
				}}
			},
		},
	})
}
