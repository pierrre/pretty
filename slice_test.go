package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Slice", []*prettytest.Case{
		{
			Name:  "Default",
			Value: []int{1, 2, 3},
		},
		{
			Name:            "Nil",
			Value:           []int(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowAddr",
			Value: []int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowCap",
			Value: []int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.ShowCap = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowIndexes",
			Value: []int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.ShowIndexes = true
			},
		},
		{
			Name:            "Empty",
			Value:           []int{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Truncated",
			Value: []int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Slice.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:            "UnknownType",
			Value:           []any{1, 2, 3},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowKnownTypes",
			Value: []int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type.ShowKnownTypes = true
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: []int{1, 2, 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Slice}
			},
			IgnoreBenchmark: true,
		},
	})
}
