package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Uintptr", []*testCase{
		{
			name:  "Default",
			value: uintptr(123),
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Uintptr}
			},
			ignoreBenchmark: true,
		},
	})
}
