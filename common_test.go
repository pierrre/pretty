package pretty_test

import (
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty"
)

func TestCommonValueWriterSetShowLen(t *testing.T) {
	vw := NewCommonWriter()
	vw.SetShowLen(true)
	assert.True(t, vw.Kind.Chan.ShowLen)
	assert.True(t, vw.Kind.Map.ShowLen)
	assert.True(t, vw.Kind.Slice.ShowLen)
	assert.True(t, vw.Kind.String.ShowLen)
	assert.True(t, vw.BytesHexDump.ShowLen)
	assert.True(t, vw.BytesableHexDump.ShowLen)
	assert.True(t, vw.Stringer.ShowLen)
}

func TestCommonValueWriterSetShowIndexes(t *testing.T) {
	vw := NewCommonWriter()
	vw.SetShowIndexes(true)
	assert.True(t, vw.Kind.Array.ShowIndexes)
	assert.True(t, vw.Kind.Slice.ShowIndexes)
	assert.True(t, vw.Kind.Chan.ShowIndexes)
}
