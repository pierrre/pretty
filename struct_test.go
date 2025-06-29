package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Struct", []*prettytest.Case{
		{
			Name: "Default",
			Value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
		},
		{
			Name:            "Empty",
			Value:           struct{}{},
			IgnoreBenchmark: true,
		},
		{
			Name: "FieldFilterExported",
			Value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Struct.FieldFilter = NewExportedStructFieldFilter()
			},
		},
		{
			Name: "SupportDisabled",
			Value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Struct}
			},
			IgnoreBenchmark: true,
		},
	})
}

type testStruct struct {
	Foo        int
	Bar        float64
	unexported int
}
