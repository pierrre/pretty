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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.Kind.Map.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "UnsortedExported",
			value: map[int]int{1: 2},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
			},
		},
		{
			name:  "UnsortedExportedShowType",
			value: map[int]int{1: 2},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.UnwrapInterface = nil
				vw.Type.ShowKnownTypes = true
			},
		},
		{
			name:  "UnsortedExportedTruncated",
			value: map[int]int{1: 2, 3: 4},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.Kind.Map.MaxLen = 1
			},
			ignoreResult: true,
		},
		{
			name:  "UnsortedExportedInterface",
			value: map[any]any{1: 2},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
			},
		},
		{
			name:  "UnsortedUnexported",
			value: testUnexported{v: map[int]int{1: 2}},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
			},
		},
		{
			name:  "UnsortedUnexportedTruncated",
			value: testUnexported{v: map[int]int{1: 2, 3: 4}},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.Kind.Map.MaxLen = 1
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "SortedExported",
			value: map[int]int{1: 2, 3: 4, 5: 6},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = true
			},
		},
		{
			name:  "SortedExportedTruncated",
			value: map[int]int{1: 2, 3: 4, 5: 6},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = true
				vw.Kind.Map.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SortedUnexported",
			value: testUnexported{v: map[int]int{1: 2, 3: 4, 5: 6}},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = true
			},
		},
		{
			name:  "UnknownType",
			value: map[any]any{"a": 1, "b": 2, "c": 3},
		},
		{
			name:  "ShowKnownTypes",
			value: map[string]int{"a": 1, "b": 2, "c": 3},
			configureWriter: func(vw *CommonWriter) {
				vw.Type.ShowKnownTypes = true
			},
		},
		{
			name:  "KeysString",
			value: map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			name:  "KeysStringShowInfos",
			value: map[string]int{"a": 1, "b": 2, "c": 3},
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.ShowKeysInfos = true
			},
		},
		{
			name:  "SupportDisabled",
			value: map[int]int{1: 2},
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Map}
			},
			ignoreBenchmark: true,
		},
	})
}
