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
			name:            "Nil",
			value:           (*testError)(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "Unexported",
			value: testUnexported{v: &testError{}},
			configure: func(vw *CommonValueWriter) {
				vw.CanInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:            "UnexportedCanInterface",
			value:           testUnexported{v: &testError{}},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: &testError{},
			configure: func(vw *CommonValueWriter) {
				vw.Error = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Error.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}
