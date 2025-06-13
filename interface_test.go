package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Interface", []*testCase{
		{
			name:  "Default",
			value: [1]any{123},
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
			},
		},
		{
			name:  "Nil",
			value: [1]any{nil},
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: [1]any{123},
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Interface}
			},
			ignoreBenchmark: true,
		},
	})
}
