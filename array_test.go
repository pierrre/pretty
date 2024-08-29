package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Array", []*testCase{
		{
			name:  "Default",
			value: [3]int{1, 2, 3},
		},
		{
			name:            "Empty",
			value:           [0]int{},
			ignoreBenchmark: true,
		},
		{
			name:            "UnknownType",
			value:           [3]any{1, 2, 3},
			ignoreBenchmark: true,
		},
		{
			name:  "ShowKnownTypes",
			value: [3]int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseArray.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}