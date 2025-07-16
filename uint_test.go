package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Uint", []*prettytest.Case{
		{
			Name:  "Default",
			Value: uint(123),
		},
		{
			Name:  "8",
			Value: uint8(123),
		},
		{
			Name:  "16",
			Value: uint16(123),
		},
		{
			Name:  "32",
			Value: uint32(123),
		},
		{
			Name:  "64",
			Value: uint64(123),
		},
		{
			Name:  "Ptr",
			Value: uintptr(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.ValueWriters[reflect.Uintptr] = vw.Kind.Uint
			},
		},
		{
			Name:  "SupportDisabled",
			Value: uint(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Uint}
			},
			IgnoreBenchmark: true,
		},
	})
}
