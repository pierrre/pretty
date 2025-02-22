package pretty_test

import (
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty"
)

func TestCommonValueWriterSetShowLen(t *testing.T) {
	vw := NewCommonValueWriter()
	vw.SetShowLen(true)
	assert.True(t, vw.Kind.BaseChan.ShowLen)
	assert.True(t, vw.Kind.BaseMap.ShowLen)
	assert.True(t, vw.Kind.BaseSlice.ShowLen)
	assert.True(t, vw.Kind.BaseString.ShowLen)
	assert.True(t, vw.BytesHexDump.ShowLen)
	assert.True(t, vw.BytesableHexDump.ShowLen)
	assert.True(t, vw.Stringer.ShowLen)
}

func TestCommonValueWriterSetShowIndexes(t *testing.T) {
	vw := NewCommonValueWriter()
	vw.SetShowIndexes(true)
	assert.True(t, vw.Kind.BaseArray.ShowIndexes)
	assert.True(t, vw.Kind.BaseSlice.ShowIndexes)
	assert.True(t, vw.Kind.BaseChan.ShowIndexes)
}
