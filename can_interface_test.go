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
				ciw := NewCanInterfaceWriter(vw)
				vw.ValueWriters = ValueWriters{ValueWriterFunc(func(st *State, v reflect.Value) bool {
					return !v.CanInterface() && ciw.WriteValue(st, v)
				})}
				vw.Recursion = nil
			},
		},
	})
}
