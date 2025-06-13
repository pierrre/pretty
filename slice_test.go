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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowCap",
			value: []int{1, 2, 3},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.ShowCap = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowIndexes",
			value: []int{1, 2, 3},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.ShowIndexes = true
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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.MaxLen = 2
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
			configureWriter: func(vw *CommonWriter) {
				vw.Type.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: []int{1, 2, 3},
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Slice}
			},
			ignoreBenchmark: true,
		},
	})
}
