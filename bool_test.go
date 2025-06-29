package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Bool", []*prettytest.Case{
		{
			Name:  "True",
			Value: true,
		},
		{
			Name:  "False",
			Value: false,
		},
		{
			Name:  "SupportDisabled",
			Value: true,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Bool}
			},
			IgnoreBenchmark: true,
		},
	})
}
