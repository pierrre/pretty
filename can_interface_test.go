package pretty_test

import (
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("CanInterface", []*testCase{
		{
			name: "Writer",
			value: func() any {
				i := 123
				return testUnexported{v: &i}
			}(),
			configureWriter: func(vw *CommonValueWriter) {
				vw.CanInterface = false
				vw.ValueWriters = ValueWriters{NewFilterValueWriter(NewCanInterfaceValueWriter(vw), func(v reflect.Value) bool {
					return !v.CanInterface()
				})}
				vw.RecursionCheck = false
			},
		},
	})
}
