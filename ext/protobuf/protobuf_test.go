package protobuf

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
	"github.com/pierrre/pretty"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	ConfigureDefault()
	pretty.DefaultCommonValueWriter.CanInterface = nil
}

func Test(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
	}{
		{
			name:  "Nil",
			value: (*wrapperspb.StringValue)(nil),
		},
		{
			name: "Unexported",
			value: func() any {
				type myType struct {
					unexported *wrapperspb.StringValue
				}
				return myType{
					unexported: wrapperspb.String("test"),
				}
			}(),
		},
		{
			name:  "String",
			value: wrapperspb.String("test"),
		},
		{
			name: "Api",
			value: &apipb.Api{
				Methods: []*apipb.Method{
					{},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertauto.Equal(t, tc.value)
		})
	}
}
