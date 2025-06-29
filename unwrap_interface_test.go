package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("UnwrapInterface", []*prettytest.Case{
		{
			Name:  "Default",
			Value: [1]any{123},
		},
		{
			Name:            "Nil",
			Value:           [1]any{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Disabled",
			Value: [1]any{123},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "DisabledNil",
			Value: [1]any{},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Writer",
			Value: [1]any{123},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
				vw.ValueWriters = ValueWriters{NewFilterWriter(NewUnwrapInterfaceWriter(vw), func(v reflect.Value) bool {
					return v.Kind() == reflect.Interface
				})}
			},
		},
		{
			Name:  "WriterNil",
			Value: [1]any{},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
				vw.ValueWriters = ValueWriters{NewFilterWriter(NewUnwrapInterfaceWriter(vw), func(v reflect.Value) bool {
					return v.Kind() == reflect.Interface
				})}
			},
		},
	})
}
