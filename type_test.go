package pretty_test

import (
	"reflect"
	"unsafe" //nolint:depguard // Required for test.

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
	"github.com/pierrre/pretty/internal/write"
)

func init() {
	prettytest.AddCasesPrefix("Type", []*prettytest.Case{
		{
			Name: "CustomBool",
			Value: func() any {
				type myBool bool
				return myBool(false)
			}(),
		},
		{
			Name: "CustomInt",
			Value: func() any {
				type myInt int
				return myInt(0)
			}(),
		},
		{
			Name: "CustomInt8",
			Value: func() any {
				type myInt8 int8
				return myInt8(0)
			}(),
		},
		{
			Name: "CustomInt16",
			Value: func() any {
				type myInt16 int16
				return myInt16(0)
			}(),
		},
		{
			Name: "CustomInt32",
			Value: func() any {
				type myInt32 int32
				return myInt32(0)
			}(),
		},
		{
			Name: "CustomInt64",
			Value: func() any {
				type myInt64 int64
				return myInt64(0)
			}(),
		},
		{
			Name: "CustomUint",
			Value: func() any {
				type myUint uint
				return myUint(0)
			}(),
		},
		{
			Name: "CustomUint8",
			Value: func() any {
				type myUint8 uint8
				return myUint8(0)
			}(),
		},
		{
			Name: "CustomUint16",
			Value: func() any {
				type myUint16 uint16
				return myUint16(0)
			}(),
		},
		{
			Name: "CustomUInt32",
			Value: func() any {
				type myUint32 uint32
				return myUint32(0)
			}(),
		},
		{
			Name: "CustomUInt64",
			Value: func() any {
				type myUint64 uint64
				return myUint64(0)
			}(),
		},
		{
			Name: "CustomUintptr",
			Value: func() any {
				type myUintptr uintptr
				return myUintptr(0)
			}(),
		},
		{
			Name: "CustomFloat32",
			Value: func() any {
				type myFloat32 float32
				return myFloat32(0)
			}(),
		},
		{
			Name: "CustomFloat64",
			Value: func() any {
				type myFloat64 float64
				return myFloat64(0)
			}(),
		},
		{
			Name: "CustomComplex64",
			Value: func() any {
				type myComplex64 complex64
				return myComplex64(0)
			}(),
		},
		{
			Name: "CustomComplex128",
			Value: func() any {
				type myComplex128 complex128
				return myComplex128(0)
			}(),
		},
		{
			Name: "CustomArray",
			Value: func() any {
				type myArray [1]int
				return myArray{}
			}(),
		},
		{
			Name: "CustomChan",
			Value: func() any {
				type myChan chan int
				return myChan(nil)
			}(),
		},
		{
			Name: "CustomFunc",
			Value: func() any {
				type myFunc func(int, string) (int, error)
				return myFunc(nil)
			}(),
		},
		{
			Name: "CustomMap",
			Value: func() any {
				type myMap map[int]int
				return myMap(nil)
			}(),
		},
		{
			Name: "CustomPointer",
			Value: func() any {
				type myPointer *int
				return myPointer(nil)
			}(),
		},
		{
			Name: "CustomSlice",
			Value: func() any {
				type mySlice []int
				return mySlice(nil)
			}(),
		},
		{
			Name: "CustomString",
			Value: func() any {
				type myString string
				return myString("")
			}(),
		},
		{
			Name: "CustomUnsafePointer",
			Value: func() any {
				type myUnsafePointer unsafe.Pointer
				return myUnsafePointer(nil)
			}(),
		},
		{
			Name: "ShowUnderlyingTypeDisabled",
			Value: func() any {
				type myBool bool
				return myBool(false)
			}(),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type.ShowUnderlyingType = false
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Disabled",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Writer",
			Value: 123,
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Type = nil
				vw.ValueWriters = ValueWriters{NewTypeWriter(vw.Kind.Int)}
			},
		},
	})
	prettytest.AddCasesPrefix("ByType", []*prettytest.Case{
		{
			Name:  "Default",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ByType[reflect.TypeFor[string]()] = ValueWriterFunc(func(st *State, v reflect.Value) bool {
					write.MustString(st.Writer, "custom")
					return true
				})
			},
		},
		{
			Name:  "NotFound",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ByType[reflect.TypeFor[int]()] = nil
			},
		},
		{
			Name:  "Support",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				cvw := ValueWriterFunc(func(st *State, v reflect.Value) bool {
					write.MustString(st.Writer, "custom")
					return true
				})
				btvw := NewByTypeWriters()
				btvw[reflect.TypeFor[string]()] = &SupportCheckerValueWriter{
					ValueWriter: cvw,
					SupportChecker: SupportCheckerFunc(func(typ reflect.Type) ValueWriter {
						return cvw
					}),
				}
				vw.Support.Checkers = []SupportChecker{
					btvw,
				}
			},
		},
	})
}
