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
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseString.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "Unquoted",
			value: "aaa\nbbb",
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseString.Quote = false
			},
		},
		{
			name:  "Truncated",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseString.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: 123,
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseString}
			},
			ignoreBenchmark: true,
		},
	})
}
