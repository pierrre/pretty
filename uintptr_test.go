package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Uintptr", []*prettytest.Case{
		{
			Name:  "Default",
			Value: uintptr(123),
		},
		{
			Name:  "SupportDisabled",
			Value: uintptr(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Uintptr}
			},
			IgnoreBenchmark: true,
		},
	})
}
