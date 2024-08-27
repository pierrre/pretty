package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Func", []*testCase{
		{
			name:  "Default",
			value: NewConfig,
		},
		{
			name:  "Nil",
			value: (func())(nil),
		},
		{
			name:  "ShowAddr",
			value: NewConfig,
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseFunc.ShowAddr = true
			},
			ignoreResult: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseFunc.WriteValue}
			},
		},
	})
}
