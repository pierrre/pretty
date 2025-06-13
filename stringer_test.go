package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Stringer", []*testCase{
		{
			name:  "Default",
			value: &testStringer{s: "test"},
		},
		{
			name:            "Nil",
			value:           (*testStringer)(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: &testStringer{s: "test"},
			configureWriter: func(vw *CommonWriter) {
				vw.Stringer.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "ReflectValue",
			value: reflect.ValueOf(123),
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{
					vw.Stringer,
					vw.ReflectValue,
				}
			},
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: &testStringer{s: "test"},
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: &testStringer{s: "test"},
			configureWriter: func(vw *CommonWriter) {
				vw.Stringer = nil
			},
			ignoreBenchmark: true,
		},
	})
}

type testStringer struct {
	s string
}

func (sr *testStringer) String() string {
	return sr.s
}
