package protobuf

import (
	"slices"
	"testing"

	"github.com/pierrre/assert/assertauto"
	"github.com/pierrre/pretty"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	pretty.DefaultCommonValueWriter.ConfigureTest()
	pretty.DefaultCommonValueWriter.ValueWriters = slices.Insert(pretty.DefaultCommonValueWriter.ValueWriters, 0, pretty.ValueWriter(&ValueWriter{
		ValueWriter: pretty.DefaultCommonValueWriter,
	}))
}

func Test(t *testing.T) {
	v := wrapperspb.String("test")
	s := pretty.String(v)
	assertauto.Equal(t, s)
}
