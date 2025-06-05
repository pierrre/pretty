package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Func", []*testCase{
		{
			name:  "Default",
			value: String,
		},
		{
			name:            "Nil",
			value:           (func())(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: String,
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.Func.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: String,
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Func}
			},
			ignoreBenchmark: true,
		},
	})
}
