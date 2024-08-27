package pretty_test

import (
	"errors"
	"io"
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Panic", []*testCase{
		{
			name:         "String",
			value:        "test",
			panicRecover: true,
			configure: func(vw *CommonValueWriter) {
				vw.PanicRecover.ShowStack = false
				vw.ValueWriters = []ValueWriter{func(w io.Writer, st State, v reflect.Value) bool {
					panic("string")
				}}
			},
		},
		{
			name:         "Error",
			value:        "test",
			panicRecover: true,
			configure: func(vw *CommonValueWriter) {
				vw.PanicRecover.ShowStack = false
				err := errors.New("error")
				vw.ValueWriters = []ValueWriter{func(w io.Writer, st State, v reflect.Value) bool {
					panic(err)
				}}
			},
		},
		{
			name:         "Other",
			value:        "test",
			panicRecover: true,
			configure: func(vw *CommonValueWriter) {
				vw.PanicRecover.ShowStack = false
				vw.ValueWriters = []ValueWriter{func(w io.Writer, st State, v reflect.Value) bool {
					panic(123)
				}}
			},
		},
		{
			name:         "ShowStack",
			value:        "test",
			panicRecover: true,
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = []ValueWriter{func(w io.Writer, st State, v reflect.Value) bool {
					panic("string")
				}}
			},
			ignoreResult: true,
			ignoreAllocs: true,
		},
		{
			name:         "Not",
			value:        "test",
			panicRecover: true,
		},
	})
}
