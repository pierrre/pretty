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
			name:  "SupportDisabled",
			value: float64(123.456),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Float}
			},
			ignoreBenchmark: true,
		},
	})
}
