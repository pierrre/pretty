package pretty_test

import (
	"runtime"
	"weak"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("WeakPointer", []*testCase{
		{
			name: "Default",
			value: func() any {
				type testWeakPointer struct {
					Pointer weak.Pointer[string]
					Value   *string
				}

				v := testWeakPointer{}
				s := "test"
				v.Value = &s
				v.Pointer = weak.Make(&s)
				return v
			}(),
		},
		{
			name: "GarbageCollected",
			value: func() weak.Pointer[[64]byte] {
				p := weak.Make(new([64]byte))
				runtime.GC()
				return p
			}(),
		},
		{
			name: "Nil",
			value: func() weak.Pointer[string] {
				return weak.Make[string](nil)
			}(),
		},
		{
			name: "Unexported",
			value: func() any {
				return testUnexported{
					v: weak.Make[string](nil),
				}
			}(),
			ignoreBenchmark: true,
		},
		{
			name: "Disabled",
			value: func() weak.Pointer[string] {
				return weak.Make[string](nil)
			}(),
			configureWriter: func(vw *CommonValueWriter) {
				vw.WeakPointer = nil
			},
			ignoreBenchmark: true,
		},
	})
}
