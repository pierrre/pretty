package pretty_test

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/pierrre/assert"
	. "github.com/pierrre/pretty"
)

func newMaxRunesTestPrinter(maxRunes int) *Printer {
	vw := NewCommonWriter()
	vw.ConfigureTest(true)
	p := NewPrinter(vw)
	p.MaxRunes = maxRunes
	return p
}

func TestMaxRunes(t *testing.T) {
	value := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	full := newMaxRunesTestPrinter(0).String(value)
	assert.StringNotContains(t, full, "<truncated>")

	t.Run("Default", func(t *testing.T) {
		vw := NewCommonWriter()
		vw.ConfigureTest(true)
		p := NewPrinter(vw)
		s := p.String(value)
		assert.Equal(t, full, s)
		assert.StringNotContains(t, s, "<truncated>")
	})

	t.Run("Zero", func(t *testing.T) {
		p := newMaxRunesTestPrinter(0)
		s := p.String(value)
		assert.Equal(t, full, s)
		assert.StringNotContains(t, s, "<truncated>")
	})

	t.Run("Limit", func(t *testing.T) {
		const limit = 20
		p := newMaxRunesTestPrinter(limit)
		s := p.String(value)
		assert.StringHasSuffix(t, s, "<truncated>")
		content := strings.TrimSuffix(s, "<truncated>")
		assert.True(t, strings.HasPrefix(full, content))
		assert.Equal(t, limit, utf8.RuneCountInString(content))
	})

	t.Run("NotReached", func(t *testing.T) {
		p := newMaxRunesTestPrinter(10000)
		s := p.String(value)
		assert.Equal(t, full, s)
		assert.StringNotContains(t, s, "<truncated>")
	})

	t.Run("Write", func(t *testing.T) {
		const limit = 20
		p := newMaxRunesTestPrinter(limit)
		var buf bytes.Buffer
		p.Write(&buf, value)
		s := buf.String()
		assert.StringHasSuffix(t, s, "<truncated>")
		content := strings.TrimSuffix(s, "<truncated>")
		assert.True(t, strings.HasPrefix(full, content))
		assert.Equal(t, limit, utf8.RuneCountInString(content))
	})
}
