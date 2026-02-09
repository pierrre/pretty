package pretty_test

import (
	"runtime"
	"weak"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("WeakPointer", []*prettytest.Case{
		{
			Name: "Default",
			Value: func() any {
				v := testWeakPointer{}
				s := "test"
				v.Value = &s
				v.Pointer = weak.Make(&s)
				return v
			}(),
		},
		{
			Name: "GarbageCollected",
			Value: func() weak.Pointer[[64]byte] {
				p := weak.Make(new([64]byte))
				runtime.GC()
				return p
			}(),
		},
		{
			Name: "Nil",
			Value: func() weak.Pointer[string] {
				return weak.Make[string](nil)
			}(),
		},
		{
			Name:            "Unexported",
			Value:           prettytest.Unexported(weak.Make[string](nil)),
			IgnoreBenchmark: true,
		},
		{
			Name: "SupportDisabled",
			Value: func() any {
				v := testWeakPointer{}
				s := "test"
				v.Value = &s
				v.Pointer = weak.Make(&s)
				return v
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name: "Disabled",
			Value: func() weak.Pointer[string] {
				return weak.Make[string](nil)
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.WeakPointer = nil
			},
			IgnoreBenchmark: true,
		},
	})
}

type testWeakPointer struct {
	Pointer weak.Pointer[string]
	Value   *string
}
