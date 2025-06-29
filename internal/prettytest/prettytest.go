// Package prettytest provides utilities for testing the pretty package.
package prettytest

import (
	"io"
	"testing"

	"github.com/pierrre/assert/assertauto"
	"github.com/pierrre/pretty"
)

// Case represents a test case.
type Case struct {
	Name             string
	Value            any
	ConfigurePrinter func(p *pretty.Printer)
	ConfigureWriter  func(vw *pretty.CommonWriter)
	IgnoreResult     bool
	IgnoreAllocs     bool
	IgnoreBenchmark  bool
}

func (tc *Case) newPrinter() *pretty.Printer {
	vw := pretty.NewCommonWriter()
	vw.ConfigureTest(true)
	if tc.ConfigureWriter != nil {
		tc.ConfigureWriter(vw)
	}
	p := pretty.NewPrinter(vw)
	if tc.ConfigurePrinter != nil {
		tc.ConfigurePrinter(p)
	}
	return p
}

var testCases []*Case

// AddCases adds test cases to the list of test cases.
func AddCases(tcs []*Case) {
	testCases = append(testCases, tcs...)
}

// AddCasesPrefix adds test cases with a prefix to the list of test cases.
func AddCasesPrefix(prefix string, tcs []*Case) {
	for _, tc := range tcs {
		tc.Name = prefix + "/" + tc.Name
	}
	AddCases(tcs)
}

// Test runs the tests.
func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			p := tc.newPrinter()
			if !tc.IgnoreResult {
				assertauto.Equal(t, tc.Value, assertauto.ValueStringer(p.String))
			}
			s := p.String(tc.Value)
			t.Log(s)
			if !tc.IgnoreAllocs {
				allocs, _ := assertauto.AllocsPerRun(t, 100, func() {
					t.Helper()
					p.Write(io.Discard, tc.Value)
				})
				t.Logf("allocs: %g", allocs)
			}
		})
	}
}

// Benchmark runs the benchmarks.
func Benchmark(b *testing.B) {
	for _, tc := range testCases {
		if !tc.IgnoreBenchmark {
			b.Run(tc.Name, func(b *testing.B) {
				p := tc.newPrinter()
				for b.Loop() {
					p.Write(io.Discard, tc.Value)
				}
			})
		}
	}
}

type unexported[T any] struct {
	v T
}

// Unexported returns an unexported value of type T.
func Unexported[T any](v T) any {
	return unexported[T]{
		v: v,
	}
}
