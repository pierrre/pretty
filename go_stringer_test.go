package pretty_test

import (
	"encoding/binary"
	"time"

	"github.com/pierrre/pretty"
	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("GoStringer", []*prettytest.Case{
		{
			Name:  "Time",
			Value: time.Time{},
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				vw.Time = nil
			},
		},
		{
			Name:  "BinaryLittleEndian",
			Value: binary.LittleEndian,
		},
		{
			Name:  "SupportDisabled",
			Value: binary.LittleEndian,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: binary.LittleEndian,
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				vw.GoStringer = nil
			},
			IgnoreBenchmark: true,
		},
	})
}
