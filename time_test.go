package pretty_test

import (
	"time"

	. "github.com/pierrre/pretty"
)

var testTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func init() {
	addTestCasesPrefix("Time", []*testCase{
		{
			name:  "Default",
			value: testTime,
		},
		{
			name:  "Zero",
			value: time.Time{},
		},
		{
			name: "Unexported",
			value: testUnexported{
				v: testTime,
			},
			configure: func(vw *CommonValueWriter) {
				vw.CanInterface = nil
			},
		},
		{
			name: "UnexportedCanInterface",
			value: testUnexported{
				v: &testTime,
			},
		},
		{
			name:  "Disabled",
			value: testTime,
			configure: func(vw *CommonValueWriter) {
				vw.Time = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.Time.WriteValue}
			},
		},
	})
}
