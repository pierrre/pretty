package pretty_test

import (
	"time"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/must"
)

var testTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func init() {
	addTestCasesPrefix("Time", []*testCase{
		{
			name:  "Default",
			value: testTime,
		},
		{
			name:            "Zero",
			value:           time.Time{},
			ignoreBenchmark: true,
		},
		{
			name:  "Location",
			value: testTime,
			configure: func(vw *CommonValueWriter) {
				var err error
				vw.Time.Location, err = time.LoadLocation("Europe/Paris")
				must.NoError(err)
			},
		},
		{
			name: "Unexported",
			value: testUnexported{
				v: testTime,
			},
			configure: func(vw *CommonValueWriter) {
				vw.CanInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name: "UnexportedCanInterface",
			value: testUnexported{
				v: &testTime,
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: testTime,
			configure: func(vw *CommonValueWriter) {
				vw.Time = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Time}
			},
			ignoreBenchmark: true,
		},
	})
}
