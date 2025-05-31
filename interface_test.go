package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Interface", []*testCase{
		{
			name:  "Default",
			value: [1]any{123},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = nil
			},
		},
		{
			name:  "Nil",
			value: [1]any{nil},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseInterface}
			},
			ignoreBenchmark: true,
		},
	})
}
