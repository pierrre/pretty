package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Complex", []*prettytest.Case{
		{
			Name:  "64",
			Value: complex64(123.456 + 789.123i),
		},
		{
			Name:  "128",
			Value: complex128(123.456 + 789.123i),
		},
		{
			Name:  "SupportDisabled",
			Value: complex128(123.456 + 789.123i),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Complex}
			},
			IgnoreBenchmark: true,
		},
	})
}
