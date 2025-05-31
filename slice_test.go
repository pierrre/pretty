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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowCap",
			value: []int{1, 2, 3},
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.ShowCap = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowIndexes",
			value: []int{1, 2, 3},
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.ShowIndexes = true
			},
		},
		{
			name:            "Empty",
			value:           []int{},
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: []int{1, 2, 3},
			configureWriter: func(vw *CommonValueWriter) {
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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Type.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.BaseSlice}
			},
			ignoreBenchmark: true,
		},
	})
}
