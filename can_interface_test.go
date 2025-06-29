package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("CanInterface", []*prettytest.Case{
		{
			Name: "Writer",
			Value: func() any {
				i := 123
				return prettytest.Unexported(&i)
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
				vw.ValueWriters = ValueWriters{NewFilterWriter(NewCanInterfaceWriter(vw), func(v reflect.Value) bool {
					return !v.CanInterface()
				})}
				vw.Recursion = nil
			},
		},
	})
}
