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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowAddr = true
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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
			},
		},
		{
			name: "ShowIndexes",
			value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.ShowIndexes = true
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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.MaxLen = 2
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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.ShowIndexes = true
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
			configureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.ShowIndexes = true
			},
		},
		{
			name: "SupportDisabled",
			value: func() chan int {
				c := make(chan int, 5)
				c <- 123
				return c
			}(),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Chan}
			},
			ignoreBenchmark: true,
		},
	})
}
