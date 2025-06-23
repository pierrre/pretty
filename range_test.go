package pretty_test

import (
	"sync"

	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Range", []*testCase{
		{
			name: "SyncMap",
			value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return m
			}(),
		},
		{
			name: "SyncUtilMap",
			value: func() any {
				m := new(syncutil.Map[string, string])
				m.Store("foo", "bar")
				return m
			}(),
		},
		{
			name:  "Nil",
			value: (*sync.Map)(nil),
		},
		{
			name:  "Empty",
			value: new(sync.Map),
		},
		{
			name: "Truncated",
			value: func() any {
				m := new(sync.Map)
				m.Store("a", "b")
				m.Store("c", "d")
				return m
			}(),
			configureWriter: func(vw *pretty.CommonWriter) {
				vw.Range.MaxLen = 1
			},
			ignoreResult: true,
		},
		{
			name: "Large",
			value: func() any {
				m := new(sync.Map)
				for i := range 100 {
					m.Store(i, i)
				}
				return m
			}(),
			ignoreResult: true,
		},
		{
			name: "Unexported",
			value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return testUnexported{
					v: m,
				}
			}(),
			configureWriter: func(vw *pretty.CommonWriter) {
				vw.CanInterface = nil
			},
			ignoreResult:    true,
			ignoreAllocs:    true,
			ignoreBenchmark: true,
		},
		{
			name: "SupportDisabled",
			value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return m
			}(),
			configureWriter: func(vw *pretty.CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name: "Disabled",
			value: func() any {
				m := new(sync.Map)
				m.Store("foo", "bar")
				return m
			}(),
			configureWriter: func(vw *pretty.CommonWriter) {
				vw.Range = nil
			},
			ignoreResult:    true,
			ignoreAllocs:    true,
			ignoreBenchmark: true,
		},
		{
			name:            "MethodNoMatch",
			value:           testRangeMethodNoMatch{},
			ignoreResult:    true,
			ignoreAllocs:    true,
			ignoreBenchmark: true,
		},
	})
}

type testRangeMethodNoMatch struct{}

func (testRangeMethodNoMatch) Range() {}
