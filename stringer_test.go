package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Stringer", []*testCase{
		{
			name:  "Default",
			value: &testStringer{s: "test"},
		},
		{
			name:            "Nil",
			value:           (*testStringer)(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: &testStringer{s: "test"},
			configure: func(vw *CommonValueWriter) {
				vw.Stringer.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Unexported",
			value: testUnexported{v: &testStringer{}},
			configure: func(vw *CommonValueWriter) {
				vw.CanInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:            "UnexportedCanInterface",
			value:           testUnexported{v: &testStringer{}},
			ignoreBenchmark: true,
		},
		{
			name:  "ReflectValue",
			value: reflect.ValueOf(123),
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{
					vw.Stringer.WriteValue,
					vw.ReflectValue.WriteValue,
				}
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: &testStringer{s: "test"},
			configure: func(vw *CommonValueWriter) {
				vw.Stringer = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Stringer.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}
