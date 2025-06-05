package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Support", []*testCase{
		{
			name:  "Empty",
			value: 123,
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support.Checkers = nil
			},
		},
	})
}
