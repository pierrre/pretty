package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("ValueWriters", []*prettytest.Case{
		{
			Name:  "Support",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support.Checkers = []SupportChecker{
					ValueWriters{
						vw.Kind.Int,
					},
				}
			},
		},
		{
			Name:  "SupportNot",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support.Checkers = []SupportChecker{
					ValueWriters{
						vw.Kind.String,
					},
				}
			},
		},
	})
}
