package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Invalid", []*testCase{
		{
			name:  "Base",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{ValueWriterFunc(func(st *State, v reflect.Value) bool {
					return vw.Kind.WriteValue(st, reflect.ValueOf(nil))
				})}
			},
		},
		{
			name:  "Nil",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{ValueWriterFunc(func(st *State, v reflect.Value) bool {
					return vw.WriteValue(st, reflect.ValueOf(nil))
				})}
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseInvalid}
			},
			ignoreBenchmark: true,
		},
	})
}
