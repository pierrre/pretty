package pretty_test

import (
	"iter"
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Recursion", []*prettytest.Case{
		{
			Name: "Pointer",
			Value: func() any {
				var v any
				v = &v
				return v
			}(),
		},
		{
			Name: "Slice",
			Value: func() []any {
				v := make([]any, 1)
				v[0] = v
				return v
			}(),
		},
		{
			Name: "Map",
			Value: func() map[int]any {
				v := make(map[int]any)
				v[0] = v
				return v
			}(),
		},
		{
			Name: "Func",
			Value: func() any {
				var f iter.Seq[any]
				f = func(yield func(any) bool) {
					yield(f)
				}
				return f
			}(),
		},
		{
			Name: "ShowAddr",
			Value: func() any {
				var v any
				v = &v
				return v
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Recursion.ShowAddr = true
			},
			IgnoreResult: true,
		},
		{
			Name:  "Disabled",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Recursion = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name: "Writer",
			Value: func() any {
				var v any
				v = &v
				return v
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				rvw := NewRecursionWriter(vw)
				rvw.ShowAddr = false
				vw.Recursion = nil
				vw.ValueWriters = ValueWriters{NewFilterWriter(rvw, func(typ reflect.Type) bool {
					return typ.Kind() == reflect.Pointer
				})}
			},
		},
		{
			Name: "DeepNested",
			Value: func() any {
				type node struct {
					Next *node
				}
				n := &node{}
				cur := n
				for range 1000 {
					cur.Next = &node{}
					cur = cur.Next
				}
				return n
			}(),
			IgnoreResult: true,
			IgnoreAllocs: true,
		},
	})
}
