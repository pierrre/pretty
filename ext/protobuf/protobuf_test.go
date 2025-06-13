package protobuf

import (
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	"github.com/pierrre/pretty"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	ConfigureDefault()
	pretty.DefaultCommonWriter.CanInterface = false
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
			name: "List",
			value: &structpb.ListValue{
				Values: []*structpb.Value{
					{
						Kind: &structpb.Value_StringValue{
							StringValue: "test",
						},
					},
					{
						Kind: &structpb.Value_NumberValue{
							NumberValue: 123,
						},
					},
				},
			},
		},
		{
			name: "Struct",
			value: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"test1": structpb.NewStringValue("test"),
					"test2": structpb.NewNumberValue(123),
				},
			},
		},
		{
			name:  "Api",
			value: &apipb.Api{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertauto.Equal(t, tc.value)
		})
	}
}

func TestSupports(t *testing.T) {
	// TODO improve tests of this package.
	vw := NewMessageWriter(nil)
	assert.Equal(t, vw.Supports(reflect.TypeFor[*wrapperspb.StringValue]()), pretty.ValueWriter(vw))
	assert.Zero(t, vw.Supports(reflect.TypeFor[string]()))
}
