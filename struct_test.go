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
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind.BaseStruct.FieldFilter = NewExportedStructFieldFilter()
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseStruct}
			},
			ignoreBenchmark: true,
		},
	})
}

type testStruct struct {
	Foo        int
	Bar        float64
	unexported int
}
