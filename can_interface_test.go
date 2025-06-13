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
			configureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
				vw.ValueWriters = ValueWriters{NewFilterWriter(NewCanInterfaceWriter(vw), func(v reflect.Value) bool {
					return !v.CanInterface()
				})}
				vw.Recursion = nil
			},
		},
	})
}
