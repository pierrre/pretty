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
			name: "ShowElems",
			value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowElems = true
			},
		},
		{
			name: "ShowElemsIndexes",
			value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowElems = true
				vw.Kind.BaseChan.ShowIndexes = true
			},
		},
		{
			name: "ShowElemsTruncated",
			value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowElems = true
				vw.Kind.BaseChan.MaxLen = 2
			},
		},
		{
			name: "ShowElemsClosed",
			value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowElems = true
			},
		},
		{
			name: "ShowElemsReadOnly",
			value: func() <-chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowElems = true
				vw.Kind.BaseChan.ShowIndexes = true
			},
		},
		{
			name: "ShowElemsWriteOnly",
			value: func() chan<- string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.Kind.BaseChan.ShowElems = true
				vw.Kind.BaseChan.ShowIndexes = true
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.BaseChan}
			},
			ignoreBenchmark: true,
		},
	})
}
