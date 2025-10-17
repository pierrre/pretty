package pretty_test

import (
	"time"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/prettytest"
)

var (
	testTimeTime     = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	testTimeDuration = 12*time.Hour + 34*time.Minute + 56*time.Second + 789*time.Millisecond
	testTimeLocation *time.Location
)

func init() {
	var err error
	testTimeLocation, err = time.LoadLocation("Europe/Paris")
	must.NoError(err)
}

func init() {
	prettytest.AddCasesPrefix("Time/Time", []*prettytest.Case{
		{
			Name:  "Default",
			Value: testTimeTime,
		},
		{
			Name:            "Zero",
			Value:           time.Time{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Location",
			Value: testTimeTime,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Time.Time.Location = testTimeLocation
			},
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(testTimeTime),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:            "UnexportedCanInterface",
			Value:           prettytest.Unexported(&testTimeTime),
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: testTimeTime,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: testTimeTime,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Time.Time = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Time.Time}
			},
			IgnoreBenchmark: true,
		},
	})
	prettytest.AddCasesPrefix("Time/Duration", []*prettytest.Case{
		{
			Name:  "Default",
			Value: testTimeDuration,
		},
		{
			Name:  "SupportDisabled",
			Value: testTimeDuration,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: testTimeDuration,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Time.Duration = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Time.Duration}
			},
			IgnoreBenchmark: true,
		},
	})
	prettytest.AddCasesPrefix("Time/Location", []*prettytest.Case{
		{
			Name:  "Default",
			Value: testTimeLocation,
		},
		{
			Name:  "Nil",
			Value: (*time.Location)(nil),
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(testTimeLocation),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:            "UnexportedCanInterface",
			Value:           prettytest.Unexported(&testTimeLocation),
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: testTimeLocation,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: testTimeLocation,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Time.Location = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.Time.Location}
			},
			IgnoreBenchmark: true,
		},
	})
}
