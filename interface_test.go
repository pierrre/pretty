package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Interface", []*prettytest.Case{
		{
			Name:  "Default",
			Value: [1]any{123},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
			},
		},
		{
			Name:  "Nil",
			Value: [1]any{nil},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: [1]any{123},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Interface}
			},
			IgnoreBenchmark: true,
		},
	})
}
