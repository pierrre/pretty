package pretty_test

import (
	"math/big"
	"strings"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("MathBigInt", []*testCase{
		{
			name:  "Default",
			value: big.NewInt(123),
		},
		{
			name: "Large",
			value: func() any {
				i := new(big.Int)
				i, _ = i.SetString(strings.Repeat("1234567890", 10), 10)
				return i
			}(),
		},
		{
			name: "Unexported",
			value: testUnexported{
				v: big.NewInt(123),
			},
			configureWriter: func(vw *CommonWriter) {
				vw.CanInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Nil",
			value: (*big.Int)(nil),
		},
		{
			name:  "SupportDisabled",
			value: big.NewInt(123),
			configureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: big.NewInt(123),
			configureWriter: func(vw *CommonWriter) {
				vw.MathBigInt = nil
			},
		},
	})
}
