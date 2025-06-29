package pretty_test

import (
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("Chan", []*prettytest.Case{
		{
			Name: "Default",
			Value: func() chan int {
				c := make(chan int, 5)
				c <- 123
				return c
			}(),
		},
		{
			Name:            "Nil",
			Value:           chan int(nil),
			IgnoreBenchmark: true,
		},
		{
			Name: "ShowAddr",
			Value: func() chan int {
				c := make(chan int, 5)
				c <- 123
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name: "ShowElems",
			Value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
			},
		},
		{
			Name: "ShowIndexes",
			Value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.ShowIndexes = true
			},
		},
		{
			Name: "ShowElemsTruncated",
			Value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.MaxLen = 2
			},
		},
		{
			Name: "ShowElemsClosed",
			Value: func() chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
			},
		},
		{
			Name: "ShowElemsReadOnly",
			Value: func() <-chan string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.ShowIndexes = true
			},
		},
		{
			Name: "ShowElemsWriteOnly",
			Value: func() chan<- string {
				c := make(chan string, 5)
				c <- "a"
				c <- "b"
				c <- "c"
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Kind.Chan.ShowElems = true
				vw.Kind.Chan.ShowIndexes = true
			},
		},
		{
			Name: "SupportDisabled",
			Value: func() chan int {
				c := make(chan int, 5)
				c <- 123
				return c
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Kind.Chan}
			},
			IgnoreBenchmark: true,
		},
	})
}
