package pretty_test

import (
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCases([]*prettytest.Case{
		{
			Name:  "Nil",
			Value: nil,
		},
	})
}
