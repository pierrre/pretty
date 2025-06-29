package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("MaxDepth", []*prettytest.Case{
		{
			Name: "Default",
			Value: func() any {
				var v1 any
				v2 := &v1
				v3 := &v2
				v4 := &v3
				return v4
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.MaxDepth.Max = 2
			},
		},
		{
			Name:  "Writer",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				mdvw := NewMaxDepthWriter(vw)
				mdvw.Max = 2
				vw.ValueWriters = ValueWriters{mdvw}
			},
		},
		{
			Name:  "WriterDisabled",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{NewMaxDepthWriter(vw.Kind.Int)}
			},
		},
	})
}
