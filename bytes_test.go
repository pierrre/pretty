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
			name:            "Nil",
			value:           []byte(nil),
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: []byte("test"),
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesHexDump.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "ShowCap",
			value: []byte("test"),
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesHexDump.ShowCap = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: []byte("test"),
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesHexDump.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: []byte("test"),
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: []byte("test"),
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesHexDump = nil
			},
			ignoreBenchmark: true,
		},
		{
			name:  "Not",
			value: "test",
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{vw.BytesHexDump}
			},
			ignoreBenchmark: true,
		},
	})
	addTestCasesPrefix("BytesableHexDump", []*testCase{
		{
			name:  "Default",
			value: &testBytesable{b: bytes.Repeat([]byte("test"), 100)},
		},
		{
			name:            "Nil",
			value:           (*testBytesable)(nil),
			ignoreBenchmark: true,
		},
		{
			name:            "NilBytes",
			value:           &testBytesable{},
			ignoreBenchmark: true,
		},
		{
			name:  "Truncated",
			value: &testBytesable{b: []byte("test")},
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesableHexDump.MaxLen = 2
			},
			ignoreBenchmark: true,
		},
		{
			name:  "ReflectValue",
			value: reflect.ValueOf(123),
			configureWriter: func(vw *CommonValueWriter) {
				vw.ValueWriters = ValueWriters{
					vw.BytesableHexDump,
					vw.ReflectValue,
				}
			},
			ignoreBenchmark: true,
		},
		{
			name:  "ShowCap",
			value: &testBytesable{b: []byte("test")},
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesableHexDump.ShowCap = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "ShowAddr",
			value: &testBytesable{b: []byte("test")},
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesableHexDump.ShowAddr = true
			},
			ignoreResult:    true,
			ignoreBenchmark: true,
		},
		{
			name:  "SupportDisabled",
			value: &testBytesable{b: []byte("test")},
			configureWriter: func(vw *CommonValueWriter) {
				vw.Support = nil
			},
		},
		{
			name:  "Disabled",
			value: &testBytesable{b: []byte("test")},
			configureWriter: func(vw *CommonValueWriter) {
				vw.BytesableHexDump = nil
			},
			ignoreBenchmark: true,
		},
	})
}

type testBytesable struct {
	b []byte
}

func (b *testBytesable) Bytes() []byte {
	return b.b
}
