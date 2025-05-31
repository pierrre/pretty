package pretty

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func Example() {
	type exampleStruct struct {
		Int    int
		Float  float64
		String string
		Map    map[string]int
		Slice  []int
	}
	v := exampleStruct{
		Int:    123,
		Float:  123.456,
		String: "test",
		Map: map[string]int{
			"foo": 1,
			"bar": 2,
		},
		Slice: []int{1, 2, 3},
	}
	s := String(v)
	fmt.Println(s)
	// Output:
	// [github.com/pierrre/pretty.exampleStruct] {
	// 	Int: [int] 123,
	// 	Float: [float64] 123.456,
	// 	String: [string] (len=4) "test",
	// 	Map: [map[string]int] (len=2) {
	// 		"bar": 2,
	// 		"foo": 1,
	// 	},
	// 	Slice: [[]int] (len=3) {
	// 		1,
	// 		2,
	// 		3,
	// 	},
	// }
}

func ExampleString() {
	s := String("test")
	fmt.Println(s)
	// Output: [string] (len=4) "test"
}

func ExampleWrite() {
	buf := new(bytes.Buffer)
	Write(buf, "test")
	s := buf.String()
	fmt.Println(s)
	// Output: [string] (len=4) "test"
}

func ExampleFormatter() {
	f := Formatter("test")
	s := fmt.Sprintf("%v", f)
	fmt.Println(s)
	// Output: [string] (len=4) "test"
}

func ExampleValueWriter() {
	vw := ValueWriterFunc(func(st *State, v reflect.Value) bool {
		_, _ = io.WriteString(st.Writer, "example")
		return true
	})
	p := NewPrinter(vw)
	s := p.String("test")
	fmt.Println(s)
	// Output: example
}
