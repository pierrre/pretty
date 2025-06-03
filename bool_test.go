package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Bool", []*testCase{
		{
			name:  "True",
			value: true,
		},
		{
			name:  "False",
			value: false,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Bool}
			},
			ignoreBenchmark: true,
		},
	})
}
