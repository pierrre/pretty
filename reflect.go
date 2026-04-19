package pretty

import (
	"reflect"
	"strconv"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/pretty/internal/itfassert"
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
		st.Writer.AppendString("<unexported>")
		return true
	}
	rv, _ := reflect.TypeAssert[reflect.Value](v)
	writeArrow(st)
	if checkInvalidNil(st, rv) {
		return true
	}
	st.KnownType = false // We want to show the type of the value.
	vw.ValueWriter.WriteValue(st, rv)
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
	st.Writer.AppendString("reflect.Type ")
	vw.writeType(st, typ)
	return true
}

func (vw *ReflectTypeWriter) writeType(st *State, typ reflect.Type) {
	st.Writer.AppendString("{\n")
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
	st.Writer.AppendByte('}')
}

func (vw *ReflectTypeWriter) writeTypeFullName(st *State, typ reflect.Type) {
	st.WriteIndent()
	st.Writer.AppendString("FullName: ")
	st.Writer.AppendString(reflectutil.TypeFullName(typ))
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypePkgPath(st *State, typ reflect.Type) {
	pkgPath := typ.PkgPath()
	if pkgPath == "" {
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("PkgPath: ")
	st.Writer.AppendString(pkgPath)
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeName(st *State, typ reflect.Type) {
	name := typ.Name()
	if name == "" {
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("Name: ")
	st.Writer.AppendString(name)
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeString(st *State, typ reflect.Type) {
	st.WriteIndent()
	st.Writer.AppendString("String: ")
	st.Writer.AppendString(typ.String())
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeKind(st *State, typ reflect.Type) {
	st.WriteIndent()
	st.Writer.AppendString("Kind: ")
	st.Writer.AppendString(typ.Kind().String())
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeSize(st *State, typ reflect.Type) {
	st.WriteIndent()
	st.Writer.AppendString("Size: ")
	st.Writer = strconv.AppendUint(st.Writer, uint64(typ.Size()), 10)
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeUnderlying(st *State, typ reflect.Type) {
	uTyp := reflectutil.GetUnderlyingType(typ)
	if uTyp == typ {
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("Underlying: ")
	vw.writeType(st, uTyp)
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeLen(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Array {
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("Len: ")
	st.Writer = strconv.AppendInt(st.Writer, int64(typ.Len()), 10)
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeKey(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Map {
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("Key: ")
	vw.writeType(st, typ.Key())
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeElem(st *State, typ reflect.Type) {
	switch typ.Kind() { //nolint:exhaustive // Kinds supported by Elem().
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
	default:
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("Elem: ")
	vw.writeType(st, typ.Elem())
	st.Writer.AppendString(",\n")
}

func (vw *ReflectTypeWriter) writeTypeChan(st *State, typ reflect.Type) {
	if typ.Kind() != reflect.Chan {
		return
	}
	st.WriteIndent()
	st.Writer.AppendString("ChanDir: ")
	st.Writer.AppendString(typ.ChanDir().String())
	st.Writer.AppendString(",\n")
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
	st.Writer.AppendString("Fields: {\n")
	st.IndentLevel++
	for _, f := range fields.Range {
		st.WriteIndent()
		st.Writer.AppendString(f.Name)
		st.Writer.AppendString(" ")
		st.Writer.AppendString(reflectutil.TypeFullName(f.Type))
		st.Writer.AppendString(",\n")
	}
	st.IndentLevel--
	st.WriteIndent()
	st.Writer.AppendString("},\n")
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
	st.Writer.AppendString(name)
	st.Writer.AppendString(": (\n")
	st.IndentLevel++
	for i := range count {
		if ignoreFirst && i == 0 {
			continue
		}
		typ := get(i)
		st.WriteIndent()
		st.Writer.AppendString(reflectutil.TypeFullName(typ))
		st.Writer.AppendString(",\n")
	}
	st.IndentLevel--
	st.WriteIndent()
	st.Writer.AppendString("),\n")
}

func (vw *ReflectTypeWriter) writeTypeMethods(st *State, typ reflect.Type) {
	methods := reflectutil.GetMethods(typ)
	if methods.Len() == 0 {
		return
	}
	ignoreReceiver := typ.Kind() != reflect.Interface
	st.WriteIndent()
	st.Writer.AppendString("Methods: {\n")
	st.IndentLevel++
	for _, m := range methods.Range {
		st.WriteIndent()
		st.Writer.AppendString(m.Name)
		st.Writer.AppendString(": {\n")
		st.IndentLevel++
		vw.writeTypeFunc(st, m.Type, ignoreReceiver)
		st.IndentLevel--
		st.WriteIndent()
		st.Writer.AppendString("},\n")
	}
	st.IndentLevel--
	st.WriteIndent()
	st.Writer.AppendString("},\n")
}

// Supports implements [SupportChecker].
func (vw *ReflectTypeWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if reflectTypeImplementsCache.ImplementedBy(typ) {
		res = vw
	}
	return res
}
