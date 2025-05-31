package pretty_test

import (
	"bytes"
	"io"
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("ReflectType", []*testCase{
		{
			name:  "String",
			value: reflect.TypeFor[string](),
		},
		{
			name:  "Pointer",
			value: reflect.TypeFor[*string](),
		},
		{
			name:  "Array",
			value: reflect.TypeFor[[10]string](),
		},
		{
			name:  "Slice",
			value: reflect.TypeFor[[]string](),
		},
		{
			name:  "Map",
			value: reflect.TypeFor[map[string]int](),
		},
		{
			name:  "Chan",
			value: reflect.TypeFor[chan<- int](),
		},
		{
			name:  "Func",
			value: reflect.TypeFor[func(int) int](),
		},
		{
			name:  "Interface",
			value: reflect.TypeFor[io.Writer](),
		},
		{
			name: "Struct",
			value: reflect.TypeFor[struct {
				String string
				Int    int
				Float  float64
				Bool   bool
			}](),
		},
		{
			name:  "EmptyInterface",
			value: reflect.TypeFor[any](),
		},
		{
			name:  "EmptyStruct",
			value: reflect.TypeFor[struct{}](),
		},
		{
			name: "CustomString",
			value: func() reflect.Type {
				type CustomString string
				return reflect.TypeFor[CustomString]()
			}(),
		},
		{
			name: "CustomPointer",
			value: func() reflect.Type {
				type CustomPointer *string
				return reflect.TypeFor[CustomPointer]()
			}(),
		},
		{
			name: "CustomSlice",
			value: func() reflect.Type {
				type CustomSlice []string
				return reflect.TypeFor[CustomSlice]()
			}(),
		},
		{
			name: "CustomStruct",
			value: func() reflect.Type {
				type CustomStruct struct {
					String string
					Int    int
					Float  float64
					Bool   bool
				}
				return reflect.TypeFor[CustomStruct]()
			}(),
		},
		{
			name:  "BytesBuffer",
			value: reflect.TypeFor[*bytes.Buffer](),
		},
		{
			name:  "Nil",
			value: [1]reflect.Type{},
			configureWriter: func(vw *CommonValueWriter) {
				vw.UnwrapInterface = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: reflect.TypeFor[string](),
			configureWriter: func(vw *CommonValueWriter) {
				vw.ReflectType = nil
			},
			ignoreBenchmark: true,
		},
	})
}
