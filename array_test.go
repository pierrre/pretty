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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.Array.ShowIndexes = true
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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.Array.MaxLen = 2
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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Type.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: [...]int{1, 2, 3},
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Array}
			},
			ignoreBenchmark: true,
		},
	})
}
