package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("ValueWriters", []*testCase{
		{
			name:  "Support",
			value: 123,
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support.Checkers = []SupportChecker{
					ValueWriters{
						&vw.Kind.Int,
					},
				}
			},
		},
		{
			name:  "SupportNot",
			value: 123,
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support.Checkers = []SupportChecker{
					ValueWriters{
						&vw.Kind.String,
					},
				}
			},
		},
	})
}
