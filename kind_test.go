package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Kind", []*prettytest.Case{
		{
			Name:  "Custom",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.ValueWriters[reflect.String] = ValueWriterFunc(func(st *State, v reflect.Value) bool {
					st.Writer.AppendString("custom")
					return true
				})
			},
		},
		{
			Name:  "Disabled",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind = nil
			},
		},
	})
}
