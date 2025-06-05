package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/write"
)

func init() {
	addTestCasesPrefix("Kind", []*testCase{
		{
			name:  "Custom",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.ValueWriters[reflect.String] = ValueWriterFunc(func(st *State, v reflect.Value) bool {
					write.MustString(st.Writer, "custom")
					return true
				})
			},
		},
	})
}
