package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Support", []*prettytest.Case{
		{
			Name:  "Empty",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support.Checkers = nil
			},
		},
	})
}
