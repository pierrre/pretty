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
			name:  "SupportDisabled",
			value: uintptr(123),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Uintptr}
			},
			ignoreBenchmark: true,
		},
	})
}
