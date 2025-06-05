package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Complex", []*testCase{
		{
			name:  "64",
			value: complex64(123.456 + 789.123i),
		},
		{
			name:  "128",
			value: complex128(123.456 + 789.123i),
		},
		{
			name:  "SupportDisabled",
			value: complex128(123.456 + 789.123i),
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Complex}
			},
			ignoreBenchmark: true,
		},
	})
}
