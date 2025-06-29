package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Invalid", []*prettytest.Case{
		{
			Name:  "Default",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{ValueWriterFunc(func(st *State, v reflect.Value) bool {
					return vw.Kind.WriteValue(st, reflect.ValueOf(nil))
				})}
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Invalid}
			},
			IgnoreBenchmark: true,
		},
	})
}
