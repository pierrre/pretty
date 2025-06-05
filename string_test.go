package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("String", []*testCase{
		{
			name:  "Default",
			value: "test",
		},
		{
			name:            "Empty",
			value:           "",
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.String.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "Unquoted",
			value: "aaa\nbbb",
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.String.Quote = false
			},
		},
		{
			name:  "Truncated",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.String.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: 123,
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.String}
			},
			ignoreBenchmark: true,
		},
	})
}
