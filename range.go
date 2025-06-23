package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// RangeWriter is a [ValueWriter] that handles values that implement the "Range" [iter.Seq2] method such as [sync.Map].
//
// It should be created with [NewRangeWriter].
type RangeWriter struct {
	ValueWriter
	// ShowKeysInfos shows keys infos.
	// Default: false.
	ShowKeysInfos bool
	// MaxLen is the maximum length of the iterator.
	// Default: 0 (no limit).
	MaxLen int
}

// NewRangeWriter creates a new [RangeWriter] with the given [IterSeq2Writer].
func NewRangeWriter(vw ValueWriter) *RangeWriter {
	return &RangeWriter{
		ValueWriter:   vw,
		ShowKeysInfos: false,
		MaxLen:        0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *RangeWriter) WriteValue(st *State, v reflect.Value) bool {
	m, ok := vw.getMethod(v.Type())
	if !ok {
		return false
	}
	if v.Kind() == reflect.Pointer && checkNil(st.Writer, v) {
		return true
	}
	if !v.CanInterface() {
		return false
	}
	writeArrowWrappedString(st.Writer, "Range() ")
	write.MustString(st.Writer, "{")
	st.IndentLevel++
	i := 0
	m.Func.Call([]reflect.Value{v, reflect.MakeFunc(m.Type.In(1), func(args []reflect.Value) []reflect.Value {
		if i == 0 {
			write.MustString(st.Writer, "\n")
		}
		st.WriteIndent()
		if vw.MaxLen > 0 && i >= vw.MaxLen {
			writeTruncated(st.Writer)
			write.MustString(st.Writer, "\n")
			return rangeReturnFalse
		}
		showInfos := st.ShowInfos
		st.ShowInfos = vw.ShowKeysInfos
		must.Handle(vw.ValueWriter.WriteValue(st, args[0]))
		st.ShowInfos = showInfos
		write.MustString(st.Writer, ": ")
		must.Handle(vw.ValueWriter.WriteValue(st, args[1]))
		write.MustString(st.Writer, ",\n")
		i++
		return rangeReturnTrue
	})})
	st.IndentLevel--
	if i != 0 {
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

// Supports implements [SupportChecker].
func (vw *RangeWriter) Supports(typ reflect.Type) ValueWriter {
	_, ok := vw.getMethod(typ)
	if ok {
		return vw
	}
	return nil
}

func (vw *RangeWriter) getMethod(typ reflect.Type) (reflect.Method, bool) {
	m, ok := reflectutil.GetMethods(typ).GetByName("Range")
	if !ok {
		return m, false
	}
	mt := m.Type
	if mt.NumIn() != 2 || mt.NumOut() != 0 {
		return m, false
	}
	ft := mt.In(1)
	return m, ft.Kind() == reflect.Func && ft.NumIn() == 2 && ft.NumOut() == 1 && ft.Out(0) == rangeBoolType
}

var (
	rangeReturnTrue  = []reflect.Value{reflect.ValueOf(true)}
	rangeReturnFalse = []reflect.Value{reflect.ValueOf(false)}
	rangeBoolType    = reflect.TypeFor[bool]()
)
