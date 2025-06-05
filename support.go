package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/syncutil"
)

// SupportChecker checks if a [reflect.Type] is supported.
// If the [reflect.Type] is supported, it returns a non nil [ValueWriter].
type SupportChecker interface {
	Supports(typ reflect.Type) ValueWriter
}

// SupportCheckerFunc is a [SupportChecker] function.
type SupportCheckerFunc func(typ reflect.Type) ValueWriter

// Supports implements [SupportChecker].
func (f SupportCheckerFunc) Supports(typ reflect.Type) ValueWriter {
	return f(typ)
}

// SupportCheckerValueWriter implements [ValueWriter] and [SupportChecker].
type SupportCheckerValueWriter struct {
	ValueWriter
	SupportChecker
}

// SupportWriter is a [ValueWriter] that selects a [ValueWriter] based on the [reflect.Type] of the [reflect.Value].
// It selects the first [SupportChecker] that supports the [reflect.Type].
//
// It should be created with [NewSupportWriter].
type SupportWriter struct {
	cache    syncutil.Map[reflect.Type, ValueWriterFunc]
	Checkers []SupportChecker
}

// NewSupportWriter creates a new [SupportWriter].
func NewSupportWriter() *SupportWriter {
	return &SupportWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *SupportWriter) WriteValue(st *State, v reflect.Value) bool {
	if len(vw.Checkers) == 0 {
		return false
	}
	typ := v.Type()
	f, ok := vw.cache.Load(typ)
	if !ok {
		for _, c := range vw.Checkers {
			w := c.Supports(typ)
			if w != nil {
				f = w.WriteValue
				break
			}
		}
		vw.cache.Store(typ, f)
	}
	if f == nil {
		return false
	}
	return f(st, v)
}

func supportsValueWriter(typ reflect.Type, vw ValueWriter) ValueWriter {
	var res ValueWriter
	c, ok := vw.(SupportChecker)
	if ok {
		res = c.Supports(typ)
	}
	return res
}
