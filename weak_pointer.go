package pretty

import (
	"reflect"
	"strings"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// WeakPointerWriter is a [ValueWriter] that handles [weak.Pointer].
//
// It should be created with [NewWeakPointerWriter].
type WeakPointerWriter struct {
	ValueWriter
}

// NewWeakPointerWriter creates a new [WeakPointerWriter] with default values.
func NewWeakPointerWriter(vw ValueWriter) *WeakPointerWriter {
	return &WeakPointerWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *WeakPointerWriter) WriteValue(st *State, v reflect.Value) bool {
	if !vw.match(v) {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	if v.IsZero() {
		writeNil(st.Writer)
		return true
	}
	m, _ := reflectutil.GetMethods(v.Type()).GetByName("Value")
	p := m.Func.Call([]reflect.Value{v})[0]
	if p.IsNil() {
		write.MustString(st.Writer, "<garbage collected>")
		return true
	}
	writeArrow(st.Writer)
	st.KnownType = false // We want to show the type of the pointer.
	must.Handle(vw.ValueWriter.WriteValue(st, p))
	return true
}

func (vw *WeakPointerWriter) match(v reflect.Value) bool {
	kind := v.Kind()
	if kind != reflect.Struct {
		return false
	}
	typ := v.Type()
	return typ.PkgPath() == "weak" && strings.HasPrefix(typ.Name(), "Pointer[")
}
