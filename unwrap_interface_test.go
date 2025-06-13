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
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
			},
			ignoreBenchmark: true,
		},
		{
			name:  "DisabledNil",
			value: [1]any{},
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Writer",
			value: [1]any{123},
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
				vw.ValueWriters = ValueWriters{NewFilterWriter(NewUnwrapInterfaceWriter(vw), func(v reflect.Value) bool {
					return v.Kind() == reflect.Interface
				})}
			},
		},
		{
			name:  "WriterNil",
			value: [1]any{},
			configureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = false
				vw.ValueWriters = ValueWriters{NewFilterWriter(NewUnwrapInterfaceWriter(vw), func(v reflect.Value) bool {
					return v.Kind() == reflect.Interface
				})}
			},
		},
	})
}
