package pretty_test

import (
	"unsafe" //nolint:depguard // Required for test.

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("UnsafePointer", []*testCase{
		{
			name: "Default",
			value: func() unsafe.Pointer {
				i := 123
				return unsafe.Pointer(&i)
			}(),
			ignoreResult: true,
		},
		{
			name: "Nil",
			value: func() unsafe.Pointer {
				return unsafe.Pointer(nil)
			}(),
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseUnsafePointer.WriteValue}
			},
		},
	})
}
