package pretty_test

import (
	"time"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/prettytest"
)

var testTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func init() {
	prettytest.AddCasesPrefix("Time", []*prettytest.Case{
		{
			Name:  "Default",
			Value: testTime,
		},
		{
			Name:            "Zero",
			Value:           time.Time{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Location",
			Value: testTime,
			ConfigureWriter: func(vw *CommonWriter) {
				var err error
				vw.Time.Location, err = time.LoadLocation("Europe/Paris")
				must.NoError(err)
			},
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(testTime),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:            "UnexportedCanInterface",
			Value:           prettytest.Unexported(&testTime),
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: testTime,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: testTime,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Time = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Time}
			},
			IgnoreBenchmark: true,
		},
	})
}
