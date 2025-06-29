package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("ReflectValue", []*prettytest.Case{
		{
			Name:  "Default",
			Value: reflect.ValueOf(123),
		},
		{
			Name:            "Nil",
			Value:           reflect.ValueOf(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "Unexported",
			Value:           prettytest.Unexported(reflect.ValueOf(123)),
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: reflect.ValueOf(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: reflect.ValueOf(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ReflectValue = nil
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.ReflectValue}
			},
			IgnoreBenchmark: true,
		},
	})
}
