package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Float", []*testCase{
		{
			name:  "32",
			value: float32(123.456),
		},
		{
			name:  "64",
			value: float64(123.456),
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.BaseFloat}
			},
			ignoreBenchmark: true,
		},
	})
}
