package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Int", []*prettytest.Case{
		{
			Name:  "Default",
			Value: 123,
		},
		{
			Name:  "8",
			Value: int8(123),
		},
		{
			Name:  "16",
			Value: int16(123),
		},
		{
			Name:  "32",
			Value: int32(123),
		},
		{
			Name:  "64",
			Value: int64(123),
		},
		{
			Name:  "SupportDisabled",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Int}
			},
			IgnoreBenchmark: true,
		},
	})
}
