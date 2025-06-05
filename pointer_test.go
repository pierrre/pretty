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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.Pointer.ShowAddr = true
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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Type.ShowKnownTypes = true
			},
			ignoreBenchmark: true,
		},
		{
			name: "SupportDisabled",
			value: func() *int {
				i := 123
				return &i
			}(),
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{&vw.Kind.Pointer}
			},
			ignoreBenchmark: true,
		},
	})
}
