package protobuf_test

import (
	"fmt"

	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/ext/protobuf"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	protobuf.ConfigureDefault()
}

func Example() {
	v := wrapperspb.String("test")
	s := pretty.String(v)
	fmt.Println(s)
	// Output:
	// [*google.golang.org/protobuf/types/known/wrapperspb.StringValue] {
	// 	value: [string] (len=4) "test",
	// }
}
