package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/write"
)

func init() {
	addTestCasesPrefix("Filter", []*testCase{
		{
			name:  "Match",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{NewFilterValueWriter(
					ValueWriterFunc(func(st *State, v reflect.Value) bool {
						write.MustString(st.Writer, "aaaa")
						return true
					}),
					func(v reflect.Value) bool {
						return v.Type() == reflect.TypeFor[string]()
					})}
			},
		},
		{
			name:  "NoMatch",
			value: 123,
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{NewFilterValueWriter(
					ValueWriterFunc(func(st *State, v reflect.Value) bool {
						panic("should not be called")
					}),
					func(v reflect.Value) bool {
						return v.Type() == reflect.TypeFor[string]()
					})}
			},
		},
	})
}
