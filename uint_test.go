package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Uint", []*testCase{
		{
			name:  "Default",
			value: uint(123),
		},
		{
			name:  "8",
			value: uint8(123),
		},
		{
			name:  "16",
			value: uint16(123),
		},
		{
			name:  "32",
			value: uint32(123),
		},
		{
			name:  "64",
			value: uint64(123),
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseUint}
			},
			ignoreBenchmark: true,
		},
	})
}
