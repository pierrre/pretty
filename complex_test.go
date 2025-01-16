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
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseComplex}
			},
			ignoreBenchmark: true,
		},
	})
}
