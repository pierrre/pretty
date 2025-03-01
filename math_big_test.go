package pretty_test

import (
	"math/big"
)

func init() {
	addTestCasesPrefix("MathBigInt", []*testCase{
		{
			name:  "Default",
			value: big.NewInt(123),
		},
	})
}
