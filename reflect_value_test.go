package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("ReflectValue", []*testCase{
		{
			name:  "Default",
			value: reflect.ValueOf(123),
		},
		{
			name:  "Nil",
			value: reflect.ValueOf(nil),
		},
		{
			name: "Unexported",
			value: testUnexported{
				v: reflect.ValueOf(123),
			},
		},
		{
			name:  "Disabled",
			value: reflect.ValueOf(123),
			configure: func(vw *CommonValueWriter) {
				vw.ReflectValue = nil
			},
			ignoreResult: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.ReflectValue.WriteValue}
			},
		},
	})
}
