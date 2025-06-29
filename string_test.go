package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("String", []*prettytest.Case{
		{
			Name:  "Default",
			Value: "test",
		},
		{
			Name:            "Empty",
			Value:           "",
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowAddr",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.String.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "Unquoted",
			Value: "aaa\nbbb",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.String.Quote = false
			},
		},
		{
			Name:  "Truncated",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.String.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.String}
			},
			IgnoreBenchmark: true,
		},
	})
}
