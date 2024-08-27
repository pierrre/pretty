package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Error", []*testCase{
		{
			name:  "Default",
			value: &testError{},
		},
		{
			name:  "Nil",
			value: (*testError)(nil),
		},
		{
			name:  "Unexported",
			value: testUnexported{v: &testError{}},
			configure: func(vw *CommonValueWriter) {
				vw.CanInterface = nil
			},
		},
		{
			name:  "UnexportedCanInterface",
			value: testUnexported{v: &testError{}},
		},
		{
			name:  "Disabled",
			value: &testError{},
			configure: func(vw *CommonValueWriter) {
				vw.Error = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Error.WriteValue}
			},
		},
	})
}
