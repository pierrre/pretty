package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Slice", []*testCase{
		{
			name:  "Default",
			value: []int{1, 2, 3},
		},
		{
			name:  "Nil",
			value: []int(nil),
		},
		{
			name:  "ShowAddr",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.ShowAddr = true
			},
			ignoreResult: true,
		},
		{
			name:  "ShowLenCapDisabled",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.SetShowLen(false)
				vw.SetShowCap(false)
			},
		},
		{
			name:  "Empty",
			value: []int{},
		},
		{
			name:  "Truncated",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseSlice.MaxLen = 2
			},
		},
		{
			name:  "UnknownType",
			value: []any{1, 2, 3},
		},
		{
			name:  "ShowKnownTypes",
			value: []int{1, 2, 3},
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowKnownTypes = true
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseSlice.WriteValue}
			},
		},
	})
}
