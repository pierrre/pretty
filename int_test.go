package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Int", []*testCase{
		{
			name:  "Default",
			value: 123,
		},
		{
			name:  "8",
			value: int8(123),
		},
		{
			name:  "16",
			value: int16(123),
		},
		{
			name:  "32",
			value: int32(123),
		},
		{
			name:  "64",
			value: int64(123),
		},
		{
			name:  "SupportDisabled",
			value: 123,
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Int}
			},
			ignoreBenchmark: true,
		},
	})
}
