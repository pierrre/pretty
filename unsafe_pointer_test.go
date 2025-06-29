package pretty_test

import (
	"unsafe" //nolint:depguard // Required for test.

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("UnsafePointer", []*prettytest.Case{
		{
			Name: "Default",
			Value: func() unsafe.Pointer {
				i := 123
				return unsafe.Pointer(&i) //nolint:gosec // It's OK.
			}(),
			IgnoreResult: true,
		},
		{
			Name: "Nil",
			Value: func() unsafe.Pointer {
				return unsafe.Pointer(nil) //nolint:gosec // It's OK.
			}(),
			IgnoreBenchmark: true,
		},
		{
			Name: "SupportDisabled",
			Value: func() unsafe.Pointer {
				i := 123
				return unsafe.Pointer(&i) //nolint:gosec // It's OK.
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
			IgnoreResult: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.UnsafePointer}
			},
			IgnoreBenchmark: true,
		},
	})
}
