package pretty_test

import (
	"errors"
	"fmt"
	"io"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
	"github.com/pierrre/pretty/internal/write"
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
			Name: "Verbose",
			Value: &testVerboseError{
				error: errors.New("error"),
			},
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

type testVerboseError struct {
	error
}

func (e *testVerboseError) Error() string {
	return "test"
}

func (e *testVerboseError) Unwrap() error {
	return e.error
}

func (e *testVerboseError) ErrorVerbose(w io.Writer) {
	write.MustString(w, "verbose a\nb\nc")
}
