package pretty_test

import (
	"math/big"
	"strings"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("MathBigInt", []*prettytest.Case{
		{
			Name:  "Default",
			Value: big.NewInt(123),
		},
		{
			Name: "Large",
			Value: func() any {
				i := new(big.Int)
				i, _ = i.SetString(strings.Repeat("1234567890", 10), 10)
				return i
			}(),
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(big.NewInt(123)),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Nil",
			Value: (*big.Int)(nil),
		},
		{
			Name:  "SupportDisabled",
			Value: big.NewInt(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: big.NewInt(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.MathBigInt = nil
			},
		},
	})
}
