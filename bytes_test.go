package pretty_test

import (
	"bytes"
	"reflect"

	. "github.com/pierrre/pretty"
	"github.com/pierrre/pretty/internal/prettytest"
)

func init() {
	prettytest.AddCasesPrefix("BytesHexDump", []*prettytest.Case{
		{
			Name:  "Default",
			Value: bytes.Repeat([]byte("test"), 100),
		},
		{
			Name:            "Nil",
			Value:           []byte(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:  "Truncated",
			Value: []byte("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesHexDump.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowCap",
			Value: []byte("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesHexDump.ShowCap = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowAddr",
			Value: []byte("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesHexDump.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: []byte("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: []byte("test"),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesHexDump = nil
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Not",
			Value: "test",
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{vw.BytesHexDump}
			},
			IgnoreBenchmark: true,
		},
	})
	prettytest.AddCasesPrefix("BytesableHexDump", []*prettytest.Case{
		{
			Name:  "Default",
			Value: &testBytesable{b: bytes.Repeat([]byte("test"), 100)},
		},
		{
			Name:            "Nil",
			Value:           (*testBytesable)(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "NilBytes",
			Value:           &testBytesable{},
			IgnoreBenchmark: true,
		},
		{
			Name:  "Truncated",
			Value: &testBytesable{b: []byte("test")},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesableHexDump.MaxLen = 2
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ReflectValue",
			Value: reflect.ValueOf(123),
			ConfigureWriter: func(vw *CommonWriter) {
				vw.ValueWriters = ValueWriters{
					vw.BytesableHexDump,
					vw.ReflectValue,
				}
			},
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowCap",
			Value: &testBytesable{b: []byte("test")},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesableHexDump.ShowCap = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "ShowAddr",
			Value: &testBytesable{b: []byte("test")},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesableHexDump.ShowAddr = true
			},
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name:  "SupportDisabled",
			Value: &testBytesable{b: []byte("test")},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.Support = nil
			},
		},
		{
			Name:  "Disabled",
			Value: &testBytesable{b: []byte("test")},
			ConfigureWriter: func(vw *CommonWriter) {
				vw.BytesableHexDump = nil
			},
			IgnoreBenchmark: true,
		},
	})
}

type testBytesable struct {
	b []byte
}

func (b *testBytesable) Bytes() []byte {
	return b.b
}
