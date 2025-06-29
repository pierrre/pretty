package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
	"github.com/pierrre/pretty/internal/write"
)

func init() {
	prettytest.AddCasesPrefix("Kind", []*prettytest.Case{
		{
			Name:  "Custom",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.ValueWriters[reflect.String] = ValueWriterFunc(func(st *State, v reflect.Value) bool {
					write.MustString(st.Writer, "custom")
					return true
				})
			},
		},
	})
}
