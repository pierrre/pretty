package protobuf

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
	"github.com/pierrre/pretty"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	ConfigureDefault()
	pretty.DefaultCommonValueWriter.CanInterface = nil
}

func Test(t *testing.T) {
	v := wrapperspb.String("test")
	assertauto.Equal(t, v)
}

func TestNil(t *testing.T) {
	var v *wrapperspb.StringValue
	assertauto.Equal(t, v)
}

func TestUnexported(t *testing.T) {
	v := struct {
		unexported *wrapperspb.StringValue
	}{
		unexported: wrapperspb.String("test"),
	}
	assertauto.Equal(t, v)
}
