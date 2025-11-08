package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Stringer", []*prettytest.Case{
		{
			Name:  "Default",
			Value: &testStringer{s: "test"},
		},
		{
			Name:            "Nil",
			Value:           (*testStringer)(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:  "Truncated",
			Value: &testStringer{s: "test"},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Stringer.ValueWriter.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ReflectValue",
			Value: reflect.ValueOf(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{
					vw.Stringer,
					vw.Reflect.Value,
				}
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Panic",
			Value: &testStringer{panic: true},
		},
		{
			Name:  "SupportDisabled",
			Value: &testStringer{s: "test"},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: &testStringer{s: "test"},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Stringer = nil
			},
			IgnoreBenchmark: true,
		},
	})
}

type testStringer struct {
	s     string
	panic bool
}

func (sr *testStringer) String() string {
	if sr.panic {
		panic("panic")
	}
	return sr.s
}
