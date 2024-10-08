package pretty_test

import (
	"unsafe" //nolint:depguard // Required for test.

	"github.com/pierrre/go-libs/reflectutil"
	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("Type", []*testCase{
		{
			name: "FullName",
			value: testStruct{
				Foo:        123,
				Bar:        123.456,
				unexported: 123,
			},
			configure: func(vw *CommonValueWriter) {
				vw.Type.Stringer = reflectutil.TypeFullName
			},
		},
		{
			name:  "Disabled",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.Type = nil
			},
			ignoreBenchmark: true,
		},
	})
	addTestCasesPrefix("TypeAndValue", []*testCase{
		{
			name: "CustomBool",
			value: func() any {
				type myBool bool
				return myBool(false)
			}(),
		},
		{
			name: "CustomInt",
			value: func() any {
				type myInt int
				return myInt(0)
			}(),
		},
		{
			name: "CustomInt8",
			value: func() any {
				type myInt8 int8
				return myInt8(0)
			}(),
		},
		{
			name: "CustomInt16",
			value: func() any {
				type myInt16 int16
				return myInt16(0)
			}(),
		},
		{
			name: "CustomInt32",
			value: func() any {
				type myInt32 int32
				return myInt32(0)
			}(),
		},
		{
			name: "CustomInt64",
			value: func() any {
				type myInt64 int64
				return myInt64(0)
			}(),
		},
		{
			name: "CustomUint",
			value: func() any {
				type myUint uint
				return myUint(0)
			}(),
		},
		{
			name: "CustomUint8",
			value: func() any {
				type myUint8 uint8
				return myUint8(0)
			}(),
		},
		{
			name: "CustomUint16",
			value: func() any {
				type myUint16 uint16
				return myUint16(0)
			}(),
		},
		{
			name: "CustomUInt32",
			value: func() any {
				type myUint32 uint32
				return myUint32(0)
			}(),
		},
		{
			name: "CustomUInt64",
			value: func() any {
				type myUint64 uint64
				return myUint64(0)
			}(),
		},
		{
			name: "CustomUintptr",
			value: func() any {
				type myUintptr uintptr
				return myUintptr(0)
			}(),
		},
		{
			name: "CustomFloat32",
			value: func() any {
				type myFloat32 float32
				return myFloat32(0)
			}(),
		},
		{
			name: "CustomFloat64",
			value: func() any {
				type myFloat64 float64
				return myFloat64(0)
			}(),
		},
		{
			name: "CustomComplex64",
			value: func() any {
				type myComplex64 complex64
				return myComplex64(0)
			}(),
		},
		{
			name: "CustomComplex128",
			value: func() any {
				type myComplex128 complex128
				return myComplex128(0)
			}(),
		},
		{
			name: "CustomArray",
			value: func() any {
				type myArray [1]int
				return myArray{}
			}(),
		},
		{
			name: "CustomChan",
			value: func() any {
				type myChan chan int
				return myChan(nil)
			}(),
		},
		{
			name: "CustomFunc",
			value: func() any {
				type myFunc func(int, string) (int, error)
				return myFunc(nil)
			}(),
		},
		{
			name: "CustomMap",
			value: func() any {
				type myMap map[int]int
				return myMap(nil)
			}(),
		},
		{
			name: "CustomPointer",
			value: func() any {
				type myPointer *int
				return myPointer(nil)
			}(),
		},
		{
			name: "CustomSlice",
			value: func() any {
				type mySlice []int
				return mySlice(nil)
			}(),
		},
		{
			name: "CustomString",
			value: func() any {
				type myString string
				return myString("")
			}(),
		},
		{
			name: "CustomUnsafePointer",
			value: func() any {
				type myUnsafePointer unsafe.Pointer
				return myUnsafePointer(nil)
			}(),
		},
		{
			name: "ShowBaseTypeDisabled",
			value: func() any {
				type myBool bool
				return myBool(false)
			}(),
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue.ShowBaseType = false
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Disabled",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.TypeAndValue = nil
			},
			ignoreBenchmark: true,
		},
	})
}
