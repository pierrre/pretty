package prettytest_test

import (
	"testing"

	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("PrettyTest", []*prettytest.Case{
		{
			Name:             "Default",
			Value:            prettytest.Unexported("test"),
			ConfigurePrinter: func(p *pretty.Printer) {},
			ConfigureWriter:  func(vw *pretty.CommonWriter) {},
		},
	})
}

func Test(t *testing.T) {
	prettytest.Test(t)
	testing.Benchmark(prettytest.Benchmark)
}
