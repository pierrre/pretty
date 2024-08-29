package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Chan", []*testCase{
		{
			name: "Default",
			value: func() chan int {
				c := make(chan int, 5)
				c <- 123
				return c
			}(),
		},
		{
			name:            "Nil",
			value:           chan int(nil),
			ignoreBenchmark: true,
		},
		{
			name: "ShowAddr",
			value: func() chan int {
				c := make(chan int, 5)
				c <- 123
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseChan.WriteValue}
			},
			ignoreBenchmark: true,
		},
	})
}
