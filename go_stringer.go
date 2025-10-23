package pretty

import (
	"fmt"
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/write"
)

var goStringerImplementsCache = reflectutil.NewImplementsCacheFor[fmt.GoStringer]()

// GoStringerWriter is a [ValueWriter] that handles [fmt.GoStringer].
//
// It should be created with [NewGoStringerWriter].
type GoStringerWriter struct{}

// NewGoStringerWriter creates a new [GoStringerWriter].
func NewGoStringerWriter() *GoStringerWriter {
	return &GoStringerWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *GoStringerWriter) WriteValue(st *State, v reflect.Value) bool {
	typ := v.Type()
	if !goStringerImplementsCache.ImplementedBy(typ) {
		return false
	}
	gsr, ok := itfassert.Assert[fmt.GoStringer](v)
	if !ok {
		return false
	}
	s := gsr.GoString()
	writeArrowWrappedString(st.Writer, "GoString() ")
	write.MustString(st.Writer, s)
	return true
}

// Supports implements [SupportChecker].
func (vw *GoStringerWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if goStringerImplementsCache.ImplementedBy(typ) {
		res = vw
	}
	return res
}
