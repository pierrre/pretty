package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
	"github.com/pierrre/pretty/internal/write"
)

func init() {
	prettytest.AddCasesPrefix("Filter", []*prettytest.Case{
		{
			Name:  "Match",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = []ValueWriter{
					NewFilterWriter(
						ValueWriterFunc(func(st *State, v reflect.Value) bool {
							write.MustString(st.Writer, "aaaa")
							return true
						}),
						func(v reflect.Value) bool {
							return v.Type() == reflect.TypeFor[string]()
						},
					),
				}
			},
		},
		{
			Name:  "NoMatch",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = []ValueWriter{
					NewFilterWriter(
						ValueWriterFunc(func(st *State, v reflect.Value) bool {
							panic("should not be called")
						}),
						func(v reflect.Value) bool {
							return v.Type() == reflect.TypeFor[string]()
						},
					),
				}
			},
		},
		{
			Name:  "Nil",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = []ValueWriter{
					NewFilterWriter(
						ValueWriterFunc(func(st *State, v reflect.Value) bool {
							write.MustString(st.Writer, "aaaa")
							return true
						}),
						nil,
					),
				}
			},
		},
	})
}
