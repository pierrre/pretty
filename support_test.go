package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Support", []*testCase{
		{
			name:  "Empty",
			value: 123,
			configureWriter: func(vw *CommonWriter) {
				vw.Support.Checkers = nil
			},
		},
	})
}
