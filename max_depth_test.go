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
			configureWriter: func(vw *CommonWriter) {
				vw.MaxDepth = 2
			},
		},
		{
			name:  "Writer",
			value: 123,
			configureWriter: func(vw *CommonWriter) {
				mdvw := NewMaxDepthWriter(vw)
				mdvw.Max = 2
				vw.ValueWriters = ValueWriters{mdvw}
			},
		},
		{
			name:  "WriterDisabled",
			value: 123,
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{NewMaxDepthWriter(vw.Kind.Int)}
			},
		},
	})
}
