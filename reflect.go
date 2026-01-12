package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// ReflectWriter is a [ValueWriter] that handles [reflect.Value] and [reflect.Type].
//
// It should be created with [NewReflectWriter].
type ReflectWriter struct {
	Value *ReflectValueWriter
	Type  *ReflectTypeWriter
}

// NewReflectWriter creates a new [ReflectWriter].
func NewReflectWriter(vw ValueWriter) *ReflectWriter {
	return &ReflectWriter{
		Value: NewReflectValueWriter(vw),
		Type:  NewReflectTypeWriter(),
	}
}

// WriteValue implements [ValueWriter].
func (vw *ReflectWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.Value != nil && vw.Value.WriteValue(st, v) {
		return true
	}
	if vw.Type != nil && vw.Type.WriteValue(st, v) {
		return true
	}
	return false
}

// Supports implements [SupportChecker].
func (vw *ReflectWriter) Supports(typ reflect.Type) ValueWriter {
	if w := callSupportCheckerPointer(vw.Value, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Type, typ); w != nil {
		return w
	}
	return nil
}

var reflectValueType = reflect.TypeFor[reflect.Value]()

// ReflectValueWriter is a [ValueWriter] that handles [reflect.Value].
//
// It should be created with [NewReflectValueWriter].
type ReflectValueWriter struct {
	ValueWriter
}

// NewReflectValueWriter creates a new [ReflectValueWriter].
func NewReflectValueWriter(vw ValueWriter) *ReflectValueWriter {
	return &ReflectValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ReflectValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct || v.Type() != reflectValueType {
		return false
	}
	if !v.CanInterface() {
		write.MustString(st.Writer, "<unexported>")
		return true
	}
	rv, _ := reflect.TypeAssert[reflect.Value](v)
	writeArrow(st.Writer)
	if checkInvalidNil(st.Writer, rv) {
		return true
	}
	st.KnownType = false // We want to show the type of the value.
	must.Handle(vw.ValueWriter.WriteValue(st, rv))
	return true
}

// Supports implements [SupportChecker].
func (vw *ReflectValueWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == reflectValueType {
		res = vw
	}
	return res
}

var reflectTypeImplementsCache = reflectutil.NewImplementsCacheFor[reflect.Type]()

// ReflectTypeWriter is a [ValueWriter] that handles [reflect.Type].
type ReflectTypeWriter struct{}

// NewReflectTypeWriter returns a new [ReflectTypeWriter].
func NewReflectTypeWriter() *ReflectTypeWriter {
	return &ReflectTypeWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *ReflectTypeWriter) WriteValue(st *State, v reflect.Value) bool {
	if !reflectTypeImplementsCache.ImplementedBy(v.Type()) {
		return false
	}
	typ, ok := itfassert.Assert[reflect.Type](v)
	if !ok {
		return false
	}
	write.MustString(st.Writer, "reflect.Type ")
	vw.writeType(st, typ)
	return true
}

func (vw *ReflectTypeWriter) writeType(st *State, typ reflect.Type) {
	write.MustString(st.Writer, "{\n")
	st.IndentLevel++
	vw.writeTypeFullName(st, typ)
	vw.writeTypePkgPath(st, typ)
	vw.writeTypeName(st, typ)
	vw.writeTypeString(st, typ)
	vw.writeTypeKind(st, typ)
	vw.writeTypeSize(st, typ)
	vw.writeTypeUnderlying(st, typ)
	vw.writeTypeLen(st, typ)
	vw.writeTypeKey(st, typ)
	vw.writeTypeElem(st, typ)
	vw.writeTypeChan(st, typ)
	vw.writeTypeStruct(st, typ)
	vw.writeTypeFunc(st, typ, false)
	vw.writeTypeMethods(st, typ)
	st.IndentLevel--
	st.WriteIndent()
	write.MustString(st.Writer, "}")
}

func (vw *ReflectTypeWriter) writeTypeFullName(st *State, typ reflect.Type) {
	st.WriteIndent()
	write.MustString(st.Writer, "FullName: ")
	write.MustString(st.Writer, reflectutil.TypeFullName(typ))
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypePkgPath(st *State, typ reflect.Type) {
	pkgPath := typ.PkgPath()
	if pkgPath == "" {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "PkgPath: ")
	write.MustString(st.Writer, pkgPath)
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeName(st *State, typ reflect.Type) {
	name := typ.Name()
	if name == "" {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "Name: ")
	write.MustString(st.Writer, name)
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeString(st *State, typ reflect.Type) {
	st.WriteIndent()
	write.MustString(st.Writer, "String: ")
	write.MustString(st.Writer, typ.String())
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeKind(st *State, typ reflect.Type) {
	st.WriteIndent()
	write.MustString(st.Writer, "Kind: ")
	write.MustString(st.Writer, typ.Kind().String())
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeSize(st *State, typ reflect.Type) {
	st.WriteIndent()
	write.MustString(st.Writer, "Size: ")
	write.Must(strconvio.WriteUint(st.Writer, uint64(typ.Size()), 10))
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeUnderlying(st *State, typ reflect.Type) {
	uTyp := reflectutil.GetUnderlyingType(typ)
	if uTyp == typ {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "Underlying: ")
	vw.writeType(st, uTyp)
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeLen(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Array {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "Len: ")
	write.Must(strconvio.WriteInt(st.Writer, int64(typ.Len()), 10))
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeKey(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Map {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "Key: ")
	vw.writeType(st, typ.Key())
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeElem(st *State, typ reflect.Type) {
	switch typ.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
	default:
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "Elem: ")
	vw.writeType(st, typ.Elem())
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeChan(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Chan {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "ChanDir: ")
	write.MustString(st.Writer, typ.ChanDir().String())
	write.MustString(st.Writer, ",\n")
}

func (vw *ReflectTypeWriter) writeTypeStruct(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Struct {
		return
	}
	fields := reflectutil.GetStructFields(typ)
	if fields.Len() == 0 {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, "Fields: {\n")
	st.IndentLevel++
	for _, f := range fields.Range {
		st.WriteIndent()
		write.MustString(st.Writer, f.Name)
		write.MustString(st.Writer, " ")
		write.MustString(st.Writer, reflectutil.TypeFullName(f.Type))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	st.WriteIndent()
	write.MustString(st.Writer, "},\n")
}

func (vw *ReflectTypeWriter) writeTypeFunc(st *State, typ reflect.Type, ignoreReceiver bool) {
	if typ.Kind() != reflect.Func {
		return
	}
	vw.writeTypeFuncParameters(st, "In", typ.NumIn(), typ.In, ignoreReceiver)
	vw.writeTypeFuncParameters(st, "Out", typ.NumOut(), typ.Out, false)
}

func (vw *ReflectTypeWriter) writeTypeFuncParameters(st *State, name string, count int, get func(int) reflect.Type, ignoreFirst bool) {
	if (ignoreFirst && count == 1) || (!ignoreFirst && count == 0) {
		return
	}
	st.WriteIndent()
	write.MustString(st.Writer, name)
	write.MustString(st.Writer, ": (\n")
	st.IndentLevel++
	for i := range count {
		if ignoreFirst && i == 0 {
			continue
		}
		typ := get(i)
		st.WriteIndent()
		write.MustString(st.Writer, reflectutil.TypeFullName(typ))
		write.MustString(st.Writer, ",\n")
	}
	st.IndentLevel--
	st.WriteIndent()
	write.MustString(st.Writer, "),\n")
}

func (vw *ReflectTypeWriter) writeTypeMethods(st *State, typ reflect.Type) {
	methods := reflectutil.GetMethods(typ)
	if methods.Len() == 0 {
		return
	}
	ignoreReceiver := typ.Kind() != reflect.Interface
	st.WriteIndent()
	write.MustString(st.Writer, "Methods: {\n")
	st.IndentLevel++
	for _, m := range methods.Range {
		st.WriteIndent()
		write.MustString(st.Writer, m.Name)
		write.MustString(st.Writer, ": {\n")
		st.IndentLevel++
		vw.writeTypeFunc(st, m.Type, ignoreReceiver)
		st.IndentLevel--
		st.WriteIndent()
		write.MustString(st.Writer, "},\n")
	}
	st.IndentLevel--
	st.WriteIndent()
	write.MustString(st.Writer, "},\n")
}

// Supports implements [SupportChecker].
func (vw *ReflectTypeWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if reflectTypeImplementsCache.ImplementedBy(typ) {
		res = vw
	}
	return res
}
