package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("MaxDepth", []*testCase{
		{
			name: "Default",
			value: func() any {
				var v1 any
				v2 := &v1
				v3 := &v2
				v4 := &v3
				return v4
			}(),
			configureWriter: func(vw *CommonValueWriter) {
				vw.MaxDepth.Max = 2
			},
		},
		{
			name: "Disabled",
			value: func() any {
				var v1 any
				v2 := &v1
				v3 := &v2
				v4 := &v3
				return v4
			}(),
			configureWriter: func(vw *CommonValueWriter) {
				vw.MaxDepth = nil
			},
			ignoreBenchmark: true,
		},
	})
}
