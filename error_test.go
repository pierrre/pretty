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
			name:  "Disabled",
			value: &testError{},
			configureWriter: func(vw *CommonValueWriter) {
				vw.Error = nil
			},
			ignoreBenchmark: true,
		},
	})
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}
