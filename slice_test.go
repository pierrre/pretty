package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Slice", []*testCase{
		{
			name:  "Default",
			value: []int{1, 2, 3},
		},
		{
			name:            "Nil",
			value:           []int(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowLenCapDisabled",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.SetShowLen(false)
				vw.SetShowCap(false)
			},
			ignoreBenchmark: true,
		},
		{
			name:            "Empty",
			value:           []int{},
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:            "UnknownType",
			value:           []any{1, 2, 3},
			ignoreBenchmark: true,
		},
		{
			name:  "ShowKnownTypes",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseSlice.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}
