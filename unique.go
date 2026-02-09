package pretty

import (
	"reflect"
	"strings"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/must"
)

// UniqueWriter is a [ValueWriter] that handles [unique.Handle].
//
// It should be created with [NewUniqueWriter].
type UniqueWriter struct {
	ValueWriter
}

// NewUniqueWriter creates a new [UniqueWriter] with default values.
func NewUniqueWriter(vw ValueWriter) *UniqueWriter {
	return &UniqueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (uw *UniqueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct {
		return false
	}
	typ := v.Type()
	if typ.PkgPath() != "unique" || !strings.HasPrefix(typ.Name(), "Handle[") {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	if v.IsZero() {
		writeNil(st.Writer)
		return true
	}
	m, _ := reflectutil.GetMethods(typ).GetByName("Value")
	v = m.Func.Call([]reflect.Value{v})[0]
	must.Handle(uw.ValueWriter.WriteValue(st, v))
	return true
}

// Supports implements [SupportChecker].
func (uw *UniqueWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Struct && typ.PkgPath() == "unique" && strings.HasPrefix(typ.Name(), "Handle[") {
		res = uw
	}
	return res
}
