package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Kind", []*testCase{
		{
			name:         "Disabled",
			value:        "test",
			panicRecover: true,
			configure: func(vw *CommonValueWriter) {
				vw.PanicRecover.ShowStack = false
				vw.Kind = nil
			},
			ignoreBenchmark: true,
		},
	})
}
