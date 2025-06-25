package pretty_test

import (
	"iter"
	"maps"
	"slices"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Iter/Seq", []*testCase{
		{
			name:  "Default",
			value: slices.Values([]string{"a", "b", "c"}),
		},
		{
			name:  "Nil",
			value: iter.Seq[int](nil),
		},
		{
			name:  "Empty",
			value: slices.Values([]string(nil)),
		},
		{
			name:  "Truncated",
			value: slices.Values([]string{"a", "b", "c"}),
			configureWriter: func(vw *CommonWriter) {
				vw.IterSeq.MaxLen = 2
			},
		},
		{
			name: "Large",
			value: func() iter.Seq[int] {
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
			name: "Unexported",
			value: testUnexported{
				v: slices.Values([]string{"a", "b", "c"}),
			},
			configureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: slices.Values([]string{"a", "b", "c"}),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: slices.Values([]string{"a", "b", "c"}),
			configureWriter: func(vw *CommonWriter) {
				vw.IterSeq = nil
			},
		},
	})
	addTestCasesPrefix("Iter/Seq2", []*testCase{
		{
			name:  "Default",
			value: slices.All([]string{"a", "b", "c"}),
		},
		{
			name:  "Nil",
			value: iter.Seq2[string, int](nil),
		},
		{
			name:  "Empty",
			value: slices.All([]string(nil)),
		},
		{
			name:  "Truncated",
			value: slices.All([]string{"a", "b", "c"}),
			configureWriter: func(vw *CommonWriter) {
				vw.IterSeq2.MaxLen = 2
			},
		},
		{
			name:  "KeysStringShowInfos",
			value: maps.All(map[string]int{"a": 1}),
			configureWriter: func(vw *CommonWriter) {
				vw.IterSeq2.ShowKeysInfos = true
			},
		},
		{
			name: "Large",
			value: func() iter.Seq2[int, int] {
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
			name: "Unexported",
			value: testUnexported{
				v: slices.All([]string{"a", "b", "c"}),
			},
			configureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: slices.All([]string{"a", "b", "c"}),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: slices.All([]string{"a", "b", "c"}),
			configureWriter: func(vw *CommonWriter) {
				vw.IterSeq2 = nil
			},
		},
	})
}
