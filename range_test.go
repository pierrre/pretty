package pretty_test

import (
	"sync"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Range", []*prettytest.Case{
		{
			Name: "SyncMap",
			Value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return m
			}(),
		},
		{
			Name: "SyncUtilMap",
			Value: func() any {
				m := new(syncutil.Map[string, string])
				m.Store("foo", "bar")
				return m
			}(),
		},
		{
			Name:  "Nil",
			Value: (*sync.Map)(nil),
		},
		{
			Name:  "Empty",
			Value: new(sync.Map),
		},
		{
			Name: "Truncated",
			Value: func() any {
				m := new(sync.Map)
				m.Store("a", "b")
				m.Store("c", "d")
				return m
			}(),
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				vw.Range.MaxLen = 1
			},
			IgnoreResult: true,
		},
		{
			Name: "Large",
			Value: func() any {
				m := new(sync.Map)
				for i := range 100 {
					m.Store(i, i)
				}
				return m
			}(),
			IgnoreResult: true,
		},
		{
			Name: "Unexported",
			Value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return prettytest.Unexported(m)
			}(),
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreResult:    true,
			IgnoreAllocs:    true,
			IgnoreBenchmark: true,
		},
		{
			Name: "SupportDisabled",
			Value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return m
			}(),
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name: "Disabled",
			Value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return m
			}(),
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				vw.Range = nil
			},
			IgnoreResult:    true,
			IgnoreAllocs:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:            "MethodNoMatch",
			Value:           testRangeMethodNoMatch{},
			IgnoreResult:    true,
			IgnoreAllocs:    true,
			IgnoreBenchmark: true,
		},
	})
}

type testRangeMethodNoMatch struct{}

func (testRangeMethodNoMatch) Range() {}
