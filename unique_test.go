package pretty_test

import (
	"unique"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Unique", []*prettytest.Case{
		{
			Name:  "Default",
			Value: unique.Make("test"),
		},
		{
			Name:  "Nil",
			Value: unique.Handle[string]{},
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(unique.Make("test")),
		},
		{
			Name:  "SupportDisabled",
			Value: unique.Make("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: unique.Make("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Unique = nil
			},
		},
	})
}
