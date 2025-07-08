package pretty_test

import (
	"math/big"
	"strings"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("MathBig/Int", []*prettytest.Case{
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
				vw.MathBig.Int = nil
			},
		},
	})
	prettytest.AddCasesPrefix("MathBig/Float", []*prettytest.Case{
		{
			Name:  "Default",
			Value: big.NewFloat(123.456),
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(big.NewFloat(123.456)),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Nil",
			Value: (*big.Float)(nil),
		},
		{
			Name:  "SupportDisabled",
			Value: big.NewFloat(123.456),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: big.NewFloat(123.456),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.MathBig.Float = nil
			},
		},
	})
	prettytest.AddCasesPrefix("MathBig/Rat", []*prettytest.Case{
		{
			Name:  "Default",
			Value: big.NewRat(355, 113),
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(big.NewRat(355, 113)),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Nil",
			Value: (*big.Rat)(nil),
		},
		{
			Name:  "SupportDisabled",
			Value: big.NewRat(355, 113),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: big.NewRat(355, 113),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.MathBig.Rat = nil
			},
		},
	})
}
