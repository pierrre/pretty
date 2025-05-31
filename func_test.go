package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Func", []*testCase{
		{
			name:  "Default",
			value: NewConfig,
		},
		{
			name:            "Nil",
			value:           (func())(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: NewConfig,
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.BaseFunc.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseFunc}
			},
			ignoreBenchmark: true,
		},
	})
}
