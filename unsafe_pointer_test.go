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
				return unsafe.Pointer(&i) //nolint:gosec // It's OK.
			}(),
			ignoreResult: true,
		},
		{
			name: "Nil",
			value: func() unsafe.Pointer {
				return unsafe.Pointer(nil) //nolint:gosec // It's OK.
			}(),
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseUnsafePointer}
			},
			ignoreBenchmark: true,
		},
	})
}
