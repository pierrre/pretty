package pretty_test

import (
	"iter"
	"maps"
	"slices"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Iter/Seq", []*prettytest.Case{
		{
			Name:  "Default",
			Value: slices.Values([]string{"a", "b", "c"}),
		},
		{
			Name:  "Nil",
			Value: iter.Seq[int](nil),
		},
		{
			Name:  "Empty",
			Value: slices.Values([]string(nil)),
		},
		{
			Name:  "Truncated",
			Value: slices.Values([]string{"a", "b", "c"}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.IterSeq.MaxLen = 2
			},
		},
		{
			Name: "Large",
			Value: func() iter.Seq[int] {
				i := 0
				return func(yield func(int) bool) {
					for {
						if i >= 100 {
							return
						}
						if !yield(i) {
							return
						}
						i++
					}
				}
			}(),
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(slices.Values([]string{"a", "b", "c"})),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: slices.Values([]string{"a", "b", "c"}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: slices.Values([]string{"a", "b", "c"}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.IterSeq = nil
			},
		},
	})
	prettytest.AddCasesPrefix("Iter/Seq2", []*prettytest.Case{
		{
			Name:  "Default",
			Value: slices.All([]string{"a", "b", "c"}),
		},
		{
			Name:  "Nil",
			Value: iter.Seq2[string, int](nil),
		},
		{
			Name:  "Empty",
			Value: slices.All([]string(nil)),
		},
		{
			Name:  "Truncated",
			Value: slices.All([]string{"a", "b", "c"}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.IterSeq2.MaxLen = 2
			},
		},
		{
			Name:  "KeysStringShowInfos",
			Value: maps.All(map[string]int{"a": 1}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.IterSeq2.ShowKeysInfos = true
			},
		},
		{
			Name: "Large",
			Value: func() iter.Seq2[int, int] {
				i := 0
				return func(yield func(int, int) bool) {
					for {
						if i >= 100 {
							return
						}
						if !yield(i, i) {
							return
						}
						i++
					}
				}
			}(),
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(slices.All([]string{"a", "b", "c"})),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: slices.All([]string{"a", "b", "c"}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: slices.All([]string{"a", "b", "c"}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.IterSeq2 = nil
			},
		},
	})
}
