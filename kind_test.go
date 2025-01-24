package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Kind", []*testCase{
		{
			name:  "Disabled",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.Kind = nil
			},
			ignoreBenchmark: true,
		},
	})
}
