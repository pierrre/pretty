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
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseInt}
			},
			ignoreBenchmark: true,
		},
	})
}
