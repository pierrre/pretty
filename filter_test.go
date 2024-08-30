package pretty_test

import (
	"io"
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal"
)

func init() {
	addTestCasesPrefix("Filter", []*testCase{
		{
			name:  "Match",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{NewFilterValueWriter(
					func(w io.Writer, st State, v reflect.Value) bool {
						internal.MustWriteString(w, "aaaa")
						return true
					},
					func(v reflect.Value) bool {
						return v.Type() == reflect.TypeFor[string]()
					}).WriteValue}
			},
		},
		{
			name:  "NoMatch",
			value: 123,
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{NewFilterValueWriter(
					func(w io.Writer, st State, v reflect.Value) bool {
						panic("should not be called")
					},
					func(v reflect.Value) bool {
						return v.Type() == reflect.TypeFor[string]()
					}).WriteValue}
			},
		},
	})
}
