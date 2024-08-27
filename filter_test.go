package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Filter", []*testCase{
		{
			name:  "Match",
			value: &testStringer{s: "test"},
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter().WriteValue, func(v reflect.Value) bool {
					return v.Type() == reflect.TypeFor[*testStringer]()
				}).WriteValue}
			},
		},
		{
			name:  "NoMatch",
			value: &testStringer{s: "test"},
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{NewFilterValueWriter(NewStringerValueWriter().WriteValue, func(v reflect.Value) bool {
					return v.Type() != reflect.TypeFor[*testStringer]()
				}).WriteValue}
			},
		},
	})
}
