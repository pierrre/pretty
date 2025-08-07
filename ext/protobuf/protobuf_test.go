package protobuf

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	ConfigureDefault()
	prettytest.AddCases([]*prettytest.Case{
		{
			Name:            "Nil",
			Value:           (*wrapperspb.StringValue)(nil),
			ConfigureWriter: ConfigureCommonWriterDefault,
		},
		{
			Name:  "Unexported",
			Value: prettytest.Unexported(wrapperspb.String("test")),
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				ConfigureCommonWriterDefault(vw)
				vw.CanInterface = nil
			},
		},
		{
			Name:            "String",
			Value:           wrapperspb.String("test"),
			ConfigureWriter: ConfigureCommonWriterDefault,
		},
		{
			Name:            "Bytes",
			Value:           wrapperspb.Bytes(bytes.Repeat([]byte("test"), 100)),
			ConfigureWriter: ConfigureCommonWriterDefault,
		},
		{
			Name: "List",
			Value: &structpb.ListValue{
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
			ConfigureWriter: ConfigureCommonWriterDefault,
		},
		{
			Name: "Struct",
			Value: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"test1": structpb.NewStringValue("test"),
					"test2": structpb.NewNumberValue(123),
				},
			},
			ConfigureWriter: ConfigureCommonWriterDefault,
		},
		{
			Name:            "Api",
			Value:           &apipb.Api{},
			ConfigureWriter: ConfigureCommonWriterDefault,
		},
		{
			Name:  "HideFieldsType",
			Value: &apipb.Api{},
			ConfigureWriter: func(vw *pretty.CommonWriter) {
				mw := NewMessageWriter(vw)
				mw.ShowFieldsType = false
				ConfigureCommonWriter(vw, mw)
			},
		},
	})
}

func Test(t *testing.T) {
	prettytest.Test(t)
}

func TestSupports(t *testing.T) {
	// TODO improve tests of this package.
	vw := NewMessageWriter(nil)
	assert.Equal(t, vw.Supports(reflect.TypeFor[*wrapperspb.StringValue]()), pretty.ValueWriter(vw))
	assert.Zero(t, vw.Supports(reflect.TypeFor[string]()))
}

func Benchmark(b *testing.B) {
	prettytest.Benchmark(b)
}
