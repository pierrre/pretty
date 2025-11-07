package pretty_test

import (
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty"
)

func TestVars(t *testing.T) {
	vs := []any{123, "abc"}
	res := Vars(vs...)
	assert.SliceEqual(t, res, vs)
}
