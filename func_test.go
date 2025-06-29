package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Func", []*prettytest.Case{
		{
			Name:  "Default",
			Value: String,
		},
		{
			Name:            "Nil",
			Value:           (func())(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowAddr",
			Value: String,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Func.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: String,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Func}
			},
			IgnoreBenchmark: true,
		},
	})
}
