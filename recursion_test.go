package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Recursion", []*testCase{
		{
			name: "Pointer",
			value: func() any {
				var v any
				v = &v
				return v
			}(),
		},
		{
			name: "Slice",
			value: func() []any {
				v := make([]any, 1)
				v[0] = v
				return v
			}(),
		},
		{
			name: "Map",
			value: func() map[int]any {
				v := make(map[int]any)
				v[0] = v
				return v
			}(),
		},
		{
			name: "ShowInfos",
			value: func() any {
				var v any
				v = &v
				return v
			}(),
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Pointer.ShowAddr = true
			},
			ignoreResult: true,
		},
		{
			name:  "Disabled",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.RecursionCheck = false
			},
			ignoreBenchmark: true,
		},
		{
			name: "Writer",
			value: func() any {
				var v any
				v = &v
				return v
			}(),
			configureWriter: func(vw *CommonWriter) {
				rvw := NewRecursionWriter(vw)
				rvw.ShowInfos = false
				vw.RecursionCheck = false
				vw.ValueWriters = ValueWriters{NewFilterWriter(rvw, func(v reflect.Value) bool {
					return v.Kind() == reflect.Pointer
				})}
			},
		},
	})
}
