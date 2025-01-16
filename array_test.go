package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Array", []*testCase{
		{
			name:  "Default",
			value: [...]int{1, 2, 3},
		},
		{
			name:  "ShowIndexes",
			value: [...]int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseArray.ShowIndexes = true
			},
		},
		{
			name:            "Empty",
			value:           [0]int{},
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: [...]int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseArray.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:            "UnknownType",
			value:           [...]any{1, 2, 3},
			ignoreBenchmark: true,
		},
		{
			name:  "ShowKnownTypes",
			value: [...]int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseArray}
			},
			ignoreBenchmark: true,
		},
	})
}
