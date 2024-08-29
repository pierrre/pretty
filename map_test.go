package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Map", []*testCase{
		{
			name:  "Nil",
			value: map[int]int(nil),
		},
		{
			name:            "Empty",
			value:           map[int]int{},
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: map[int]int{1: 2},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = false
				vw.Kind.BaseMap.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "UnsortedExported",
			value: map[int]int{1: 2},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = false
			},
		},
		{
			name:  "UnsortedExportedTruncated",
			value: map[int]int{1: 2, 3: 4},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = false
				vw.Kind.BaseMap.MaxLen = 1
			},
			ignoreResult: true,
		},
		{
			name:  "UnsortedExportedInterface",
			value: map[any]any{1: 2},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = false
			},
		},
		{
			name:  "UnsortedUnexported",
			value: testUnexported{v: map[int]int{1: 2}},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = false
			},
		},
		{
			name:  "UnsortedUnexportedTruncated",
			value: testUnexported{v: map[int]int{1: 2, 3: 4}},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = false
				vw.Kind.BaseMap.MaxLen = 1
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "SortedBool",
			value: map[bool]int{false: 1, true: 2},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
			},
		},
		{
			name:  "SortedInt",
			value: map[int]int{1: 2, 3: 4, 5: 6},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
			},
		},
		{
			name:  "SortedUint",
			value: map[uint]int{1: 2, 3: 4, 5: 6},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
			},
		},
		{
			name:  "SortedFloat",
			value: map[float64]int{1: 2, 3: 4, 5: 6},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
			},
		},
		{
			name:  "SortedString",
			value: map[string]int{"a": 1, "b": 2, "c": 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
			},
		},
		{
			name:  "SortedDefault",
			value: map[testComparableStruct]int{{V: 1}: 2, {V: 3}: 4, {V: 5}: 6},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
			},
		},
		{
			name:  "SortedDefaultSimple",
			value: map[testComparableStruct]int{{V: 1}: 2, {V: 3}: 4, {V: 5}: 6},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
				vw.Kind.BaseMap.SortKeysCmpDefault = nil
			},
			ignoreAllocs: true,
		},
		{
			name:  "SortedTruncated",
			value: map[int]int{1: 2, 3: 4, 5: 6},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseMap.SortKeys = true
				vw.Kind.BaseMap.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "UnknownType",
			value: map[any]any{"a": 1, "b": 2, "c": 3},
		},
		{
			name:  "ShowKnownTypes",
			value: map[string]int{"a": 1, "b": 2, "c": 3},
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowKnownTypes = true
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseMap.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}
