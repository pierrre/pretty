package pretty_test

import (
	"errors"
	"fmt"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Error", []*prettytest.Case{
		{
			Name: "Default",
			Value: fmt.Errorf("test: %w", errors.Join(
				errors.New("error1"),
				errors.New("error2"),
			)),
		},
		{
			Name:            "Nil",
			Value:           (*testError)(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: &testError{},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: &testError{},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Error = nil
			},
			IgnoreBenchmark: true,
		},
	})
}

type testError struct{}

func (e *testError) Error() string {
	return "test"
}
