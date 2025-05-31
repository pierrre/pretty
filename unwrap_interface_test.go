package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("UnwrapInterface", []*testCase{
		{
			name:  "Default",
			value: [1]any{123},
		},
		{
			name:            "Nil",
			value:           [1]any{},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: [1]any{123},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = false
			},
			ignoreBenchmark: true,
		},
		{
			name:  "DisabledNil",
			value: [1]any{},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = false
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Writer",
			value: [1]any{123},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = false
				vw.ValueWriters = ValueWriters{NewFilterValueWriter(NewUnwrapInterfaceValueWriter(&vw.Kind.BaseInt), func(v reflect.Value) bool {
					return v.Kind() == reflect.Interface
				})}
			},
		},
		{
			name:  "WriterNil",
			value: [1]any{},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = false
				vw.ValueWriters = ValueWriters{NewFilterValueWriter(NewUnwrapInterfaceValueWriter(nil), func(v reflect.Value) bool {
					return v.Kind() == reflect.Interface
				})}
			},
		},
	})
}
