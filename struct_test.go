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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Struct.FieldFilter = NewExportedStructFieldFilter()
			},
		},
		{
			name: "SupportDisabled",
			value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Struct}
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
