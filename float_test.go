package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Float", []*prettytest.Case{
		{
			Name:  "32",
			Value: float32(123.456),
		},
		{
			Name:  "64",
			Value: float64(123.456),
		},
		{
			Name:  "SupportDisabled",
			Value: float64(123.456),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Float}
			},
			IgnoreBenchmark: true,
		},
	})
}
