package pretty

import (
	"io"
	"reflect"
	"unsafe" //nolint:depguard // Required for unsafe.Pointer.

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// TypeValueWriter is a [ValueWriter] that writes the type of the value.
//
// It should be created with [NewTypeValueWriter].
type TypeValueWriter struct {
	ValueWriter
	// Stringer converts the [reflect.Type] to a string.
	// Default: [reflectutil.TypeFullName].
	Stringer func(reflect.Type) string
	// ShowKnownTypes shows known types.
	// Default: false.
	ShowKnownTypes bool
	// ShowBaseType shows the base type.
	// Default: true.
	ShowBaseType bool
}

// NewTypeValueWriter creates a new [TypeValueWriter] with default values.
func NewTypeValueWriter(vw ValueWriter) *TypeValueWriter {
	return &TypeValueWriter{
		ValueWriter:    vw,
		Stringer:       reflectutil.TypeFullName,
		ShowKnownTypes: false,
		ShowBaseType:   true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if !st.KnownType || vw.ShowKnownTypes {
		write.MustString(st.Writer, "[")
		write.MustString(st.Writer, vw.Stringer(v.Type()))
		write.MustString(st.Writer, "]")
		vw.writeBaseType(st.Writer, v)
		write.MustString(st.Writer, " ")
	}
	knownType := st.KnownType
	st.KnownType = true // The type is known, because we showed it.
	must.Handle(vw.ValueWriter.WriteValue(st, v))
	st.KnownType = knownType
	return true
}

func (vw *TypeValueWriter) writeType(w io.Writer, typ reflect.Type) {
	write.MustString(w, vw.Stringer(typ))
}

func (vw *TypeValueWriter) writeBaseType(w io.Writer, v reflect.Value) {
	if !vw.ShowBaseType {
		return
	}
	typ := v.Type()
	if typ.PkgPath() == "" {
		return
	}
	baseType := vw.getCachedBaseType(typ)
	if baseType == nil {
		return
	}
	write.MustString(w, "(")
	vw.writeType(w, baseType)
	write.MustString(w, ")")
}

var baseTypeCache syncutil.Map[reflect.Type, reflect.Type]

func (vw *TypeValueWriter) getCachedBaseType(typ reflect.Type) reflect.Type {
	baseType, ok := baseTypeCache.Load(typ)
	if ok {
		return baseType
	}
	baseType = vw.getBaseType(typ)
	baseTypeCache.Store(typ, baseType)
	return baseType
}

//nolint:gocyclo // We need to handle all kinds.
func (vw *TypeValueWriter) getBaseType(typ reflect.Type) reflect.Type {
	var baseType reflect.Type
	switch typ.Kind() { //nolint:exhaustive //Some kinds are not handled.
	case reflect.Bool:
		baseType = reflect.TypeFor[bool]()
	case reflect.Int:
		baseType = reflect.TypeFor[int]()
	case reflect.Int8:
		baseType = reflect.TypeFor[int8]()
	case reflect.Int16:
		baseType = reflect.TypeFor[int16]()
	case reflect.Int32:
		baseType = reflect.TypeFor[int32]()
	case reflect.Int64:
		baseType = reflect.TypeFor[int64]()
	case reflect.Uint:
		baseType = reflect.TypeFor[uint]()
	case reflect.Uint8:
		baseType = reflect.TypeFor[uint8]()
	case reflect.Uint16:
		baseType = reflect.TypeFor[uint16]()
	case reflect.Uint32:
		baseType = reflect.TypeFor[uint32]()
	case reflect.Uint64:
		baseType = reflect.TypeFor[uint64]()
	case reflect.Uintptr:
		baseType = reflect.TypeFor[uintptr]()
	case reflect.Float32:
		baseType = reflect.TypeFor[float32]()
	case reflect.Float64:
		baseType = reflect.TypeFor[float64]()
	case reflect.Complex64:
		baseType = reflect.TypeFor[complex64]()
	case reflect.Complex128:
		baseType = reflect.TypeFor[complex128]()
	case reflect.Array:
		baseType = reflect.ArrayOf(typ.Len(), typ.Elem())
	case reflect.Chan:
		baseType = reflect.ChanOf(typ.ChanDir(), typ.Elem())
	case reflect.Func:
		in := make([]reflect.Type, typ.NumIn())
		for i := range in {
			in[i] = typ.In(i)
		}
		out := make([]reflect.Type, typ.NumOut())
		for i := range out {
			out[i] = typ.Out(i)
		}
		baseType = reflect.FuncOf(in, out, typ.IsVariadic())
	case reflect.Map:
		baseType = reflect.MapOf(typ.Key(), typ.Elem())
	case reflect.Pointer:
		baseType = reflect.PointerTo(typ.Elem())
	case reflect.Slice:
		baseType = reflect.SliceOf(typ.Elem())
	case reflect.String:
		baseType = reflect.TypeFor[string]()
	case reflect.UnsafePointer:
		baseType = reflect.TypeFor[unsafe.Pointer]()
	}
	if baseType == nil {
		// Unimplemented: invalid, interface, struct.
		return nil
	}
	if baseType == typ {
		return nil
	}
	return baseType
}

// ByTypeValueWriters is a [ValueWriter] that selects a [ValueWriter] by type name.
//
// It should be created with [NewByTypeValueWriters].
type ByTypeValueWriters map[string]ValueWriter

// NewByTypeValueWriters creates a new [ByTypeValueWriters].
func NewByTypeValueWriters() ByTypeValueWriters {
	return make(ByTypeValueWriters)
}

// WriteValue implements [ValueWriter].
func (vw ByTypeValueWriters) WriteValue(st *State, v reflect.Value) bool {
	if len(vw) == 0 {
		return false
	}
	typ := v.Type()
	name := reflectutil.TypeFullName(typ)
	w, ok := vw[name]
	if !ok {
		return false
	}
	return w.WriteValue(st, v)
}
