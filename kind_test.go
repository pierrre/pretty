package pretty_test

import (
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Kind", []*testCase{
		{
			name:  "Disabled",
			value: "test",
			configurePrinter: func(p *Printer) {
				p.PanicRecover = true
			},
			configureWriter: func(vw *CommonValueWriter) {
				vw.Kind = nil
			},
			ignoreBenchmark: true,
		},
	})
}
