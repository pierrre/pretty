package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Array", []*prettytest.Case{
		{
			Name:  "Default",
			Value: [...]int{1, 2, 3},
		},
		{
			Name:  "ShowIndexes",
			Value: [...]int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Array.ShowIndexes = true
			},
		},
		{
			Name:            "Empty",
			Value:           [0]int{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Truncated",
			Value: [...]int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Array.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:            "UnknownType",
			Value:           [...]any{1, 2, 3},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowKnownTypes",
			Value: [...]int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type.ShowKnownTypes = true
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: [...]int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Array}
			},
			IgnoreBenchmark: true,
		},
	})
}
