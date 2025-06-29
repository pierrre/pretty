package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Map", []*prettytest.Case{
		{
			Name:  "Nil",
			Value: map[int]int(nil),
		},
		{
			Name:            "Empty",
			Value:           map[int]int{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowAddr",
			Value: map[int]int{1: 2},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.Kind.Map.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "UnsortedExported",
			Value: map[int]int{1: 2},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
			},
		},
		{
			Name:  "UnsortedExportedShowType",
			Value: map[int]int{1: 2},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.UnwrapInterface = nil
				vw.Type.ShowKnownTypes = true
			},
		},
		{
			Name:  "UnsortedExportedTruncated",
			Value: map[int]int{1: 2, 3: 4},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.Kind.Map.MaxLen = 1
			},
			IgnoreResult: true,
		},
		{
			Name:  "UnsortedExportedInterface",
			Value: map[any]any{1: 2},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
			},
		},
		{
			Name:  "UnsortedUnexported",
			Value: prettytest.Unexported(map[int]int{1: 2}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
			},
		},
		{
			Name:  "UnsortedUnexportedTruncated",
			Value: prettytest.Unexported(map[int]int{1: 2, 3: 4}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = false
				vw.Kind.Map.MaxLen = 1
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "SortedExported",
			Value: map[int]int{1: 2, 3: 4, 5: 6},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = true
			},
		},
		{
			Name:  "SortedExportedTruncated",
			Value: map[int]int{1: 2, 3: 4, 5: 6},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = true
				vw.Kind.Map.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SortedUnexported",
			Value: prettytest.Unexported(map[int]int{1: 2, 3: 4, 5: 6}),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.SortKeys = true
			},
		},
		{
			Name:  "UnknownType",
			Value: map[any]any{"a": 1, "b": 2, "c": 3},
		},
		{
			Name:  "ShowKnownTypes",
			Value: map[string]int{"a": 1, "b": 2, "c": 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type.ShowKnownTypes = true
			},
		},
		{
			Name:  "KeysString",
			Value: map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			Name:  "KeysStringShowInfos",
			Value: map[string]int{"a": 1, "b": 2, "c": 3},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Map.ShowKeysInfos = true
			},
		},
		{
			Name:  "SupportDisabled",
			Value: map[int]int{1: 2},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Map}
			},
			IgnoreBenchmark: true,
		},
	})
}
