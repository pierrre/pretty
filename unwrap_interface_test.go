package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("UnwrapInterface", []*testCase{
		{
			name:  "Default",
			value: [1]any{123},
		},
		{
			name:            "Nil",
			value:           [1]any{},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: [1]any{123},
			configure: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "DisabledNil",
			value: [1]any{},
			configure: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = nil
			},
			ignoreBenchmark: true,
		},
	})
}
