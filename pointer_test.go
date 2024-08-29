package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Pointer", []*testCase{
		{
			name: "Default",
			value: func() *int {
				i := 123
				return &i
			}(),
		},
		{
			name:            "Nil",
			value:           (*int)(nil),
			ignoreBenchmark: true,
		},
		{
			name: "ShowAddr",
			value: func() *int {
				i := 123
				return &i
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BasePointer.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name: "UnknownType",
			value: func() *any {
				i := any(123)
				return &i
			}(),
			ignoreBenchmark: true,
		},
		{
			name: "ShowKnownTypes",
			value: func() *int {
				i := 123
				return &i
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BasePointer.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}