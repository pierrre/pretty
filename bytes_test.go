package pretty_test

import (
	"bytes"
	"reflect"

	. "github.com/pierrre/pretty"
)

func init() {
	addTestCasesPrefix("BytesHexDump", []*testCase{
		{
			name:  "Default",
			value: bytes.Repeat([]byte("test"), 100),
		},
		{
			name:  "Nil",
			value: []byte(nil),
		},
		{
			name:  "Truncated",
			value: []byte("test"),
			configure: func(vw *CommonValueWriter) {
				vw.BytesHexDump.MaxLen = 2
			},
		},
		{
			name:  "ShowAddr",
			value: []byte("test"),
			configure: func(vw *CommonValueWriter) {
				vw.BytesHexDump.ShowAddr = true
			},
			ignoreResult: true,
		},
		{
			name:  "Disabled",
			value: []byte("test"),
			configure: func(vw *CommonValueWriter) {
				vw.BytesHexDump = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.BytesHexDump.WriteValue}
			},
		},
	})
	addTestCasesPrefix("BytesableHexDump", []*testCase{
		{
			name:  "Default",
			value: &testBytesable{b: bytes.Repeat([]byte("test"), 100)},
		},
		{
			name:  "Nil",
			value: (*testBytesable)(nil),
		},
		{
			name:  "NilBytes",
			value: &testBytesable{},
		},
		{
			name:  "Truncated",
			value: &testBytesable{b: []byte("test")},
			configure: func(vw *CommonValueWriter) {
				vw.BytesableHexDump.MaxLen = 2
			},
		},
		{
			name: "Unexported",
			value: testUnexported{
				v: &testBytesable{},
			},
			configure: func(vw *CommonValueWriter) {
				vw.CanInterface = nil
			},
		},
		{
			name: "UnexportedCanInterface",
			value: testUnexported{
				v: &testBytesable{},
			},
		},
		{
			name:  "ReflectValue",
			value: reflect.ValueOf(123),
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{
					vw.BytesableHexDump.WriteValue,
					vw.ReflectValue.WriteValue,
				}
			},
		},
		{
			name:  "ShowAddr",
			value: &testBytesable{b: []byte("test")},
			configure: func(vw *CommonValueWriter) {
				vw.BytesableHexDump.ShowAddr = true
			},
			ignoreResult: true,
		},
		{
			name:  "Disabled",
			value: &testBytesable{b: []byte("test")},
			configure: func(vw *CommonValueWriter) {
				vw.BytesableHexDump = nil
			},
		},
		{
			name:  "Not",
			value: "test",
			configure: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.BytesableHexDump.WriteValue}
			},
		},
	})
}
