package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Pointer", []*prettytest.Case{
		{
			Name: "Default",
			Value: func() *int {
				i := 123
				return &i
			}(),
		},
		{
			Name:            "Nil",
			Value:           (*int)(nil),
			IgnoreBenchmark: true,
		},
		{
			Name: "ShowAddr",
			Value: func() *int {
				i := 123
				return &i
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Pointer.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name: "UnknownType",
			Value: func() *any {
				i := any(123)
				return &i
			}(),
			IgnoreBenchmark: true,
		},
		{
			Name: "ShowKnownTypes",
			Value: func() *int {
				i := 123
				return &i
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type.ShowKnownTypes = true
			},
			IgnoreBenchmark: true,
		},
		{
			Name: "SupportDisabled",
			Value: func() *int {
				i := 123
				return &i
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Pointer}
			},
			IgnoreBenchmark: true,
		},
	})
}
