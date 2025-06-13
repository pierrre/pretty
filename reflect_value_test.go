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
			name:            "Nil",
			value:           reflect.ValueOf(nil),
			ignoreBenchmark: true,
		},
		{
			name: "Unexported",
			value: testUnexported{
				v: reflect.ValueOf(123),
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: reflect.ValueOf(123),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: reflect.ValueOf(123),
			configureWriter: func(vw *CommonWriter) {
				vw.ReflectValue = nil
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.ReflectValue}
			},
			ignoreBenchmark: true,
		},
	})
}
