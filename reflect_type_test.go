package pretty_test

import (
	"bytes"
	"io"
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("ReflectType", []*prettytest.Case{
		{
			Name:  "String",
			Value: reflect.TypeFor[string](),
		},
		{
			Name:  "Pointer",
			Value: reflect.TypeFor[*string](),
		},
		{
			Name:  "Array",
			Value: reflect.TypeFor[[10]string](),
		},
		{
			Name:  "Slice",
			Value: reflect.TypeFor[[]string](),
		},
		{
			Name:  "Map",
			Value: reflect.TypeFor[map[string]int](),
		},
		{
			Name:  "Chan",
			Value: reflect.TypeFor[chan<- int](),
		},
		{
			Name:  "Func",
			Value: reflect.TypeFor[func(int) int](),
		},
		{
			Name:  "Interface",
			Value: reflect.TypeFor[io.Writer](),
		},
		{
			Name: "Struct",
			Value: reflect.TypeFor[struct {
				String string
				Int    int
				Float  float64
				Bool   bool
			}](),
		},
		{
			Name:  "EmptyInterface",
			Value: reflect.TypeFor[any](),
		},
		{
			Name:  "EmptyStruct",
			Value: reflect.TypeFor[struct{}](),
		},
		{
			Name: "CustomString",
			Value: func() reflect.Type {
				type CustomString string
				return reflect.TypeFor[CustomString]()
			}(),
		},
		{
			Name: "CustomPointer",
			Value: func() reflect.Type {
				type CustomPointer *string
				return reflect.TypeFor[CustomPointer]()
			}(),
		},
		{
			Name: "CustomSlice",
			Value: func() reflect.Type {
				type CustomSlice []string
				return reflect.TypeFor[CustomSlice]()
			}(),
		},
		{
			Name: "CustomStruct",
			Value: func() reflect.Type {
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
			Name:  "BytesBuffer",
			Value: reflect.TypeFor[*bytes.Buffer](),
		},
		{
			Name:  "Nil",
			Value: [1]reflect.Type{},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.UnwrapInterface = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: reflect.TypeFor[string](),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: reflect.TypeFor[string](),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ReflectType = nil
			},
			IgnoreBenchmark: true,
		},
	})
}
