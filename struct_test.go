package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Struct", []*testCase{
		{
			name: "Default",
			value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
		},
		{
			name:            "Empty",
			value:           struct{}{},
			ignoreBenchmark: true,
		},
		{
			name: "FieldFilterExported",
			value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseStruct.FieldFilter = ExportedStructFieldFilter
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseStruct.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}
