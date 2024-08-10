// Package pretty provides utilities to pretty print values.
package pretty

import (
	"bytes"
	"cmp"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"sync"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/strconvio"
)

// Write writes the value to the [io.Writer] with [DefaultPrinter].
func Write(w io.Writer, vi any) {
	DefaultPrinter.Write(w, vi)
}

// String returns the value as a string with [DefaultPrinter].
func String(vi any) string {
	return DefaultPrinter.String(vi)
}

// Formatter returns a [fmt.Formatter] for the value with [DefaultPrinter].
func Formatter(vi any) fmt.Formatter {
	return DefaultPrinter.Formatter(vi)
}

// DefaultPrinter is the default [Printer] created with [NewPrinterCommon], [DefaultConfig] and [DefaultCommonValueWriter].
var DefaultPrinter = NewPrinterCommon(DefaultConfig, DefaultCommonValueWriter)

// Printer pretty-prints values.
//
// It should be created with [NewPrinter].
type Printer struct {
	Config      *Config
	ValueWriter ValueWriter
}

// NewPrinter creates a new [Printer].
func NewPrinter(c *Config, vw ValueWriter) *Printer {
	return &Printer{
		Config:      c,
		ValueWriter: vw,
	}
}

// NewPrinterCommon creates a new [Printer] with a [CommonValueWriter].
//
// It calls [CommonValueWriter.ConfigureWithPrinter] with the created [Printer].
func NewPrinterCommon(c *Config, vw *CommonValueWriter) *Printer {
	p := NewPrinter(c, vw.WriteValue)
	vw.ConfigureWithPrinter(p)
	return p
}

// Write writes the value to the [io.Writer].
func (p *Printer) Write(w io.Writer, vi any) {
	v := reflect.ValueOf(vi)
	if !v.IsValid() {
		writeNil(w)
		return
	}
	stItf := statePool.Get()
	defer statePool.Put(stItf)
	st := stItf.(*State) //nolint:forcetypeassert // The pool only contains *State.
	st.reset()
	mustHandle(p.ValueWriter(p.Config, w, st, v))
}

var bufPool = &bufpool.Pool{}

// String returns the value as a string.
func (p *Printer) String(vi any) string {
	buf := p.getBuf(vi)
	defer bufPool.Put(buf)
	return buf.String()
}

func (p *Printer) getBuf(vi any) *bytes.Buffer {
	buf := bufPool.Get()
	p.Write(buf, vi)
	return buf
}

func (p *Printer) compare(a, b reflect.Value) int {
	aBuf := p.getBuf(a)
	defer bufPool.Put(aBuf)
	bBuf := p.getBuf(b)
	defer bufPool.Put(bBuf)
	return bytes.Compare(aBuf.Bytes(), bBuf.Bytes())
}

// Formatter returns a [fmt.Formatter] for the value.
func (p *Printer) Formatter(vi any) fmt.Formatter {
	return &formatter{
		printer: p,
		value:   vi,
	}
}

// DefaultConfig is the default [Config] created with [NewConfig].
var DefaultConfig = NewConfig()

// Config is a configuration used to pretty print values.
//
// It should be created with [NewConfig].
type Config struct {
	// Indent is the string used to indent.
	// Default: "\t".
	Indent string
}

// NewConfig creates a new [Config] initialized with default values.
func NewConfig() *Config {
	return &Config{
		Indent: "\t",
	}
}

// WriteIndent writes the indentation to the [io.Writer].
func (c *Config) WriteIndent(w io.Writer, st *State) {
	WriteIndent(w, c.Indent, st.Indent)
}

var statePool = &sync.Pool{
	New: func() any {
		return &State{}
	},
}

// State represents the state of the [Printer].
//
// Functions must restore the original state when they return.
type State struct {
	Depth   int
	Indent  int
	Visited []uintptr
}

// RunDepth runs the function with incremented depth and restores the original depth after.
func (st *State) RunDepth(f func(st *State)) {
	st.Depth++
	defer func() {
		st.Depth--
	}()
	f(st)
}

// RunIndent runs the function with incremented indentation and restores the original indentation after.
func (st *State) RunIndent(f func(st *State)) {
	st.Indent++
	defer func() {
		st.Indent--
	}()
	f(st)
}

// RunVisited runs the function with the visited pointer and restores the original visited pointers after.
func (st *State) RunVisited(p uintptr, f func(st *State)) {
	l := len(st.Visited)
	st.Visited = append(st.Visited, p)
	defer func() {
		st.Visited = st.Visited[:l]
	}()
	f(st)
}

func (st *State) reset() {
	st.Depth = 0
	st.Indent = 0
	st.Visited = st.Visited[:0]
}

// ValueWriter is a function that writes a [reflect.Value] to a [io.Writer].
//
// It returns true if it handles the value, false otherwise.
// If it returns false, it must not write anything.
//
// Implementations must check [reflect.Value.CanInterface] before using [reflect.Value.Interface].
//
// Implentations can assume that the value is valid.
type ValueWriter func(c *Config, w io.Writer, st *State, v reflect.Value) bool

// KindValueWriter is a [ValueWriter] that writes the value with the kind-specific [ValueWriter].
//
// The "Base*" fields are the default [ValueWriter] for each group of kinds.
// They can be configured.
//
// The other fields are the [ValueWriter] for each kind.
// It's allowed to update theses fields.
// Each [ValueWriter] must be able to handle the corresponding kind.
//
// It should be created with [NewKindValueWriter].
type KindValueWriter struct {
	BaseInvalid       *InvalidValueWriter
	BaseBool          *BoolValueWriter
	BaseInt           *IntValueWriter
	BaseUint          *UintValueWriter
	BaseUintptr       *UintptrValueWriter
	BaseFloat         *FloatValueWriter
	BaseComplex       *ComplexValueWriter
	BaseArray         *ArrayValueWriter
	BaseChan          *ChanValueWriter
	BaseFunc          *FuncValueWriter
	BaseInterface     *InterfaceValueWriter
	BaseMap           *MapValueWriter
	BasePointer       *PointerValueWriter
	BaseSlice         *SliceValueWriter
	BaseString        *StringValueWriter
	BaseStruct        *StructValueWriter
	BaseUnsafePointer *UnsafePointerValueWriter

	Invalid       ValueWriter
	Bool          ValueWriter
	Int           ValueWriter
	Int8          ValueWriter
	Int16         ValueWriter
	Int32         ValueWriter
	Int64         ValueWriter
	Uint          ValueWriter
	Uint8         ValueWriter
	Uint16        ValueWriter
	Uint32        ValueWriter
	Uint64        ValueWriter
	Uintptr       ValueWriter
	Float32       ValueWriter
	Float64       ValueWriter
	Complex64     ValueWriter
	Complex128    ValueWriter
	Array         ValueWriter
	Chan          ValueWriter
	Func          ValueWriter
	Interface     ValueWriter
	Map           ValueWriter
	Pointer       ValueWriter
	Slice         ValueWriter
	String        ValueWriter
	Struct        ValueWriter
	UnsafePointer ValueWriter
}

// NewKindValueWriter creates a new [KindValueWriter] with default values.
func NewKindValueWriter(vw ValueWriter) *KindValueWriter {
	kindVW := &KindValueWriter{
		BaseInvalid:       NewInvalidValueWriter(),
		BaseBool:          NewBoolValueWriter(),
		BaseInt:           NewIntValueWriter(),
		BaseUint:          NewUintValueWriter(),
		BaseUintptr:       NewUintptrValueWriter(),
		BaseFloat:         NewFloatValueWriter(),
		BaseComplex:       NewComplexValueWriter(),
		BaseArray:         NewArrayValueWriter(vw),
		BaseChan:          NewChanValueWriter(),
		BaseFunc:          NewFuncValueWriter(),
		BaseInterface:     NewInterfaceValueWriter(vw),
		BaseMap:           NewMapValueWriter(vw),
		BasePointer:       NewPointerValueWriter(vw),
		BaseSlice:         NewSliceValueWriter(vw),
		BaseString:        NewStringValueWriter(),
		BaseStruct:        NewStructValueWriter(vw),
		BaseUnsafePointer: NewUnsafePointerValueWriter(),
	}
	kindVW.Invalid = kindVW.writeInvalid
	kindVW.Bool = kindVW.writeBool
	kindVW.Int = kindVW.writeInt
	kindVW.Int8 = kindVW.writeInt
	kindVW.Int16 = kindVW.writeInt
	kindVW.Int32 = kindVW.writeInt
	kindVW.Int64 = kindVW.writeInt
	kindVW.Uint = kindVW.writeUint
	kindVW.Uint8 = kindVW.writeUint
	kindVW.Uint16 = kindVW.writeUint
	kindVW.Uint32 = kindVW.writeUint
	kindVW.Uint64 = kindVW.writeUint
	kindVW.Uintptr = kindVW.writeUintptr
	kindVW.Float32 = kindVW.writeFloat
	kindVW.Float64 = kindVW.writeFloat
	kindVW.Complex64 = kindVW.writeComplex
	kindVW.Complex128 = kindVW.writeComplex
	kindVW.Array = kindVW.writeArray
	kindVW.Chan = kindVW.writeChan
	kindVW.Func = kindVW.writeFunc
	kindVW.Interface = kindVW.writeInterface
	kindVW.Map = kindVW.writeMap
	kindVW.Pointer = kindVW.writePointer
	kindVW.Slice = kindVW.writeSlice
	kindVW.String = kindVW.writeString
	kindVW.Struct = kindVW.writeStruct
	kindVW.UnsafePointer = kindVW.writeUnsafePointer
	return kindVW
}

// WriteValue implements [ValueWriter].
//
//nolint:gocyclo // We need to handle all kinds.
func (vw *KindValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // All kinds are handled, "invalid" is in the default case.
	case reflect.Bool:
		mustHandle(vw.Bool(c, w, st, v))
	case reflect.Int:
		mustHandle(vw.Int(c, w, st, v))
	case reflect.Int8:
		mustHandle(vw.Int8(c, w, st, v))
	case reflect.Int16:
		mustHandle(vw.Int16(c, w, st, v))
	case reflect.Int32:
		mustHandle(vw.Int32(c, w, st, v))
	case reflect.Int64:
		mustHandle(vw.Int64(c, w, st, v))
	case reflect.Uint:
		mustHandle(vw.Uint(c, w, st, v))
	case reflect.Uint8:
		mustHandle(vw.Uint8(c, w, st, v))
	case reflect.Uint16:
		mustHandle(vw.Uint16(c, w, st, v))
	case reflect.Uint32:
		mustHandle(vw.Uint32(c, w, st, v))
	case reflect.Uint64:
		mustHandle(vw.Uint64(c, w, st, v))
	case reflect.Uintptr:
		mustHandle(vw.Uintptr(c, w, st, v))
	case reflect.Float32:
		mustHandle(vw.Float32(c, w, st, v))
	case reflect.Float64:
		mustHandle(vw.Float64(c, w, st, v))
	case reflect.Complex64:
		mustHandle(vw.Complex64(c, w, st, v))
	case reflect.Complex128:
		mustHandle(vw.Complex128(c, w, st, v))
	case reflect.Array:
		mustHandle(vw.Array(c, w, st, v))
	case reflect.Chan:
		mustHandle(vw.Chan(c, w, st, v))
	case reflect.Func:
		mustHandle(vw.Func(c, w, st, v))
	case reflect.Interface:
		mustHandle(vw.Interface(c, w, st, v))
	case reflect.Map:
		mustHandle(vw.Map(c, w, st, v))
	case reflect.Pointer:
		mustHandle(vw.Pointer(c, w, st, v))
	case reflect.Slice:
		mustHandle(vw.Slice(c, w, st, v))
	case reflect.String:
		mustHandle(vw.String(c, w, st, v))
	case reflect.Struct:
		mustHandle(vw.Struct(c, w, st, v))
	case reflect.UnsafePointer:
		mustHandle(vw.UnsafePointer(c, w, st, v))
	default:
		mustHandle(vw.Invalid(c, w, st, v))
	}
	return true
}

func (vw *KindValueWriter) writeInvalid(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseInvalid.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeBool(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseBool.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeInt(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseInt.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeUint(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseUint.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeUintptr(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseUintptr.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeFloat(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseFloat.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeComplex(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseComplex.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeArray(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseArray.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeChan(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseChan.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeFunc(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseFunc.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeInterface(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseInterface.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeMap(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseMap.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writePointer(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BasePointer.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeSlice(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseSlice.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeString(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseString.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeStruct(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseStruct.WriteValue(c, w, st, v)
}

func (vw *KindValueWriter) writeUnsafePointer(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.BaseUnsafePointer.WriteValue(c, w, st, v)
}

// InvalidValueWriter is a [ValueWriter] that handles invalid values.
//
// It should be created with [NewInvalidValueWriter].
type InvalidValueWriter struct{}

// NewInvalidValueWriter creates a new [InvalidValueWriter].
func NewInvalidValueWriter() *InvalidValueWriter {
	return &InvalidValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *InvalidValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return checkInvalid(w, v)
}

func checkInvalid(w io.Writer, v reflect.Value) bool {
	if v.IsValid() {
		return false
	}
	writeString(w, "<invalid>")
	return true
}

// BoolValueWriter is a [ValueWriter] that handles bool values.
//
// It should be created with [NewBoolValueWriter].
type BoolValueWriter struct{}

// NewBoolValueWriter creates a new [BoolValueWriter].
func NewBoolValueWriter() *BoolValueWriter {
	return &BoolValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *BoolValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Bool {
		return false
	}
	mustWrite(strconvio.WriteBool(w, v.Bool()))
	return true
}

// IntValueWriter is a [ValueWriter] that handles int values.
//
// It should be created with [NewIntValueWriter].
type IntValueWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewIntValueWriter creates a new [IntValueWriter] with default values.
func NewIntValueWriter() *IntValueWriter {
	return &IntValueWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *IntValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		mustWrite(strconvio.WriteInt(w, v.Int(), vw.Base))
		return true
	}
	return false
}

// UintValueWriter is a [ValueWriter] that handles uint values.
//
// It should be created with [NewUintValueWriter].
type UintValueWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewUintValueWriter creates a new [UintValueWriter] with default values.
func NewUintValueWriter() *UintValueWriter {
	return &UintValueWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *UintValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles uint.
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mustWrite(strconvio.WriteUint(w, v.Uint(), vw.Base))
		return true
	}
	return false
}

// UintptrValueWriter is a [ValueWriter] that handles uintptr values.
//
// It should be created with [NewUintptrValueWriter].
type UintptrValueWriter struct{}

// NewUintptrValueWriter creates a new [UintptrValueWriter].
func NewUintptrValueWriter() *UintptrValueWriter {
	return &UintptrValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *UintptrValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Uintptr {
		return false
	}
	writeUintptr(w, v.Uint())
	return true
}

func writeUintptr(w io.Writer, p uint64) {
	writeString(w, "0x")
	mustWrite(strconvio.WriteUint(w, p, 16))
}

// FloatValueWriter is a [ValueWriter] that handles float values.
//
// It should be created with [NewFloatValueWriter].
type FloatValueWriter struct {
	// Format is the format used to format the float.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the float.
	// Default: -1.
	Precision int
}

// NewFloatValueWriter creates a new [FloatValueWriter] with default values.
func NewFloatValueWriter() *FloatValueWriter {
	return &FloatValueWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FloatValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles float.
	case reflect.Float32, reflect.Float64:
		mustWrite(strconvio.WriteFloat(w, v.Float(), vw.Format, vw.Precision, v.Type().Bits()))
		return true
	}
	return false
}

// ComplexValueWriter is a [ValueWriter] that handles complex values.
//
// It should be created with [NewComplexValueWriter].
type ComplexValueWriter struct {
	// Format is the format used to format the complex.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the complex.
	// Default: -1.
	Precision int
}

// NewComplexValueWriter creates a new [ComplexValueWriter] with default values.
func NewComplexValueWriter() *ComplexValueWriter {
	return &ComplexValueWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ComplexValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles complex.
	case reflect.Complex64, reflect.Complex128:
		writeString(w, strconv.FormatComplex(v.Complex(), vw.Format, vw.Precision, v.Type().Bits()))
		return true
	}
	return false
}

// ArrayValueWriter is a [ValueWriter] that handles array values.
//
// It should be created with [NewArrayValueWriter].
type ArrayValueWriter struct {
	ValueWriter
	// MaxLen is the maximum length of the array.
	// Default: 0 (no limit).
	MaxLen int
}

// NewArrayValueWriter creates a new [ArrayValueWriter] with default values.
func NewArrayValueWriter(vw ValueWriter) *ArrayValueWriter {
	return &ArrayValueWriter{
		ValueWriter: vw,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ArrayValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Array {
		return false
	}
	writeArray(c, w, st, v, vw.MaxLen, vw.ValueWriter)
	return true
}

func writeArray(c *Config, w io.Writer, st *State, v reflect.Value, maxLen int, vw ValueWriter) {
	l := v.Len()
	truncated := false
	if maxLen > 0 && l > maxLen {
		l = maxLen
		truncated = true
	}
	writeString(w, "{")
	if v.Len() > 0 {
		writeString(w, "\n")
		st.RunIndent(func(st *State) {
			for i := range l {
				c.WriteIndent(w, st)
				mustHandle(vw(c, w, st, v.Index(i)))
				writeString(w, ",\n")
			}
			if truncated {
				c.WriteIndent(w, st)
				writeTruncated(w)
				writeString(w, "\n")
			}
		})
		c.WriteIndent(w, st)
	}
	writeString(w, "}")
}

// ChanValueWriter is a [ValueWriter] that handles chan values.
//
// It should be created with [NewChanValueWriter].
type ChanValueWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
}

// NewChanValueWriter creates a new [ChanValueWriter] with default values.
func NewChanValueWriter() *ChanValueWriter {
	return &ChanValueWriter{
		ShowLen:  true,
		ShowCap:  true,
		ShowAddr: true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *ChanValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Chan {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	infos{
		showLen:  vw.ShowLen,
		len:      v.Len(),
		showCap:  vw.ShowCap,
		cap:      v.Cap(),
		showAddr: vw.ShowAddr,
		addr:     uint64(uintptr(v.UnsafePointer())),
	}.write(w)
	return true
}

// StructValueWriter is a [ValueWriter] that handles struct values.
//
// It should be created with [NewStructValueWriter].
type FuncValueWriter struct {
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
}

// NewFuncValueWriter creates a new [FuncValueWriter] with default values.
func NewFuncValueWriter() *FuncValueWriter {
	return &FuncValueWriter{
		ShowAddr: true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FuncValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Func {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	p := uintptr(v.UnsafePointer())
	infos{
		showAddr: vw.ShowAddr,
		addr:     uint64(p),
	}.writeWithTrailingSpace(w)
	name := runtime.FuncForPC(p).Name()
	writeString(w, name)
	return true
}

// InterfaceValueWriter is a [ValueWriter] that handles interface values.
//
// It should be created with [NewInterfaceValueWriter].
type InterfaceValueWriter struct {
	ValueWriter
}

// NewInterfaceValueWriter creates a new [InterfaceValueWriter].
func NewInterfaceValueWriter(vw ValueWriter) *InterfaceValueWriter {
	return &InterfaceValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *InterfaceValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Interface {
		return false
	}
	writeArrow(w)
	if checkNil(w, v) {
		return true
	}
	mustHandle(vw.ValueWriter(c, w, st, v.Elem()))
	return true
}

// MapValueWriter is a [ValueWriter] that handles map values.
//
// It should be created with [NewMapValueWriter].
type MapValueWriter struct {
	ValueWriter
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
	// SortKeys sorts the keys.
	// Default: false.
	SortKeys bool
	// SortKeysCmpDefault is the default comparison function for sorting the keys, when the key type is not ordered.
	// Default: a function that uses [fmt.Sprint].
	SortKeysCmpDefault func(a, b reflect.Value) int
	// MaxLen is the maximum length of the map.
	// Default: 0 (no limit).
	MaxLen int
}

// NewMapValueWriter creates a new [MapValueWriter] with default values.
func NewMapValueWriter(vw ValueWriter) *MapValueWriter {
	return &MapValueWriter{
		ValueWriter:        vw,
		ShowLen:            true,
		ShowAddr:           true,
		SortKeys:           false,
		SortKeysCmpDefault: mapSortKeysCmpDefault,
		MaxLen:             0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MapValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Map {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	infos{
		showLen:  vw.ShowLen,
		len:      v.Len(),
		showAddr: vw.ShowAddr,
		addr:     uint64(uintptr(v.UnsafePointer())),
	}.writeWithTrailingSpace(w)
	writeString(w, "{")
	if v.Len() > 0 {
		writeString(w, "\n")
		st.RunIndent(func(st *State) {
			if vw.SortKeys {
				vw.writeSorted(c, w, st, v)
			} else {
				vw.writeUnsorted(c, w, st, v)
			}
		})
		c.WriteIndent(w, st)
	}
	writeString(w, "}")
	return true
}

func (vw *MapValueWriter) writeSorted(c *Config, w io.Writer, st *State, v reflect.Value) {
	keys := vw.getSortedKeys(v)
	for i, key := range keys {
		ok := vw.writeEntry(c, w, st, key, v.MapIndex(key), i)
		if !ok {
			break
		}
	}
}

func (vw *MapValueWriter) getSortedKeys(v reflect.Value) []reflect.Value {
	keys := v.MapKeys()
	vw.sortKeys(v.Type().Key(), keys)
	return keys
}

func (vw *MapValueWriter) sortKeys(typ reflect.Type, vs []reflect.Value) {
	cmpFunc := vw.getSortKeysCmp(typ)
	slices.SortFunc(vs, cmpFunc)
}

func (vw *MapValueWriter) getSortKeysCmp(typ reflect.Type) func(a, b reflect.Value) int {
	switch typ.Kind() { //nolint:exhaustive // Optimized for common kinds, the default case is less optimized.
	case reflect.Bool:
		return func(a, b reflect.Value) int {
			if !a.Bool() {
				return -1
			}
			return 1
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.Int(), b.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.Uint(), b.Uint())
		}
	case reflect.Float32, reflect.Float64:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.Float(), b.Float())
		}
	case reflect.String:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.String(), b.String())
		}
	}
	if vw.SortKeysCmpDefault != nil {
		return vw.SortKeysCmpDefault
	}
	return mapSortKeysCmpDefault
}

func mapSortKeysCmpDefault(a, b reflect.Value) int {
	return cmp.Compare(fmt.Sprint(a), fmt.Sprint(b))
}

func (vw *MapValueWriter) writeUnsorted(c *Config, w io.Writer, st *State, v reflect.Value) {
	if v.CanInterface() {
		vw.writeUnsortedExported(c, w, st, v)
	} else {
		vw.writeUnsortedUnexported(c, w, st, v)
	}
}

var typeInterface = reflect.TypeFor[any]()

var reflectValuePool = &sync.Pool{
	New: func() any {
		return reflect.New(typeInterface).Elem()
	},
}

func (vw *MapValueWriter) writeUnsortedExported(c *Config, w io.Writer, st *State, v reflect.Value) {
	iter := v.MapRange()
	keyItf := reflectValuePool.Get()
	valueItf := reflectValuePool.Get()
	key := keyItf.(reflect.Value)     //nolint:forcetypeassert // The pool only contains reflect.Value.
	value := valueItf.(reflect.Value) //nolint:forcetypeassert // The pool only contains reflect.Value.
	defer func() {
		key.SetZero()
		value.SetZero()
		reflectValuePool.Put(keyItf)
		reflectValuePool.Put(valueItf)
	}()
	for i := 0; iter.Next(); i++ {
		key.SetIterKey(iter)
		value.SetIterValue(iter)
		ok := vw.writeEntry(c, w, st, key, value, i)
		if !ok {
			break
		}
	}
}

func (vw *MapValueWriter) writeUnsortedUnexported(c *Config, w io.Writer, st *State, v reflect.Value) {
	iter := v.MapRange()
	for i := 0; iter.Next(); i++ {
		key := iter.Key()
		value := iter.Value()
		ok := vw.writeEntry(c, w, st, key, value, i)
		if !ok {
			break
		}
	}
}

func (vw *MapValueWriter) writeEntry(c *Config, w io.Writer, st *State, key reflect.Value, value reflect.Value, i int) bool {
	c.WriteIndent(w, st)
	if vw.MaxLen > 0 && i >= vw.MaxLen {
		writeTruncated(w)
		writeString(w, "\n")
		return false
	}
	mustHandle(vw.ValueWriter(c, w, st, key))
	writeString(w, ": ")
	mustHandle(vw.ValueWriter(c, w, st, value))
	writeString(w, ",\n")
	return true
}

// PointerValueWriter is a [ValueWriter] that handles pointer values.
//
// It should be created with [NewPointerValueWriter].
type PointerValueWriter struct {
	ValueWriter
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
}

// NewPointerValueWriter creates a new [PointerValueWriter] with default values.
func NewPointerValueWriter(vw ValueWriter) *PointerValueWriter {
	return &PointerValueWriter{
		ValueWriter: vw,
		ShowAddr:    true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *PointerValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	infos{
		showAddr: vw.ShowAddr,
		addr:     uint64(uintptr(v.UnsafePointer())),
	}.writeWithTrailingSpace(w)
	writeArrow(w)
	mustHandle(vw.ValueWriter(c, w, st, v.Elem()))
	return true
}

// SliceValueWriter is a [ValueWriter] that handles slice values.
//
// It should be created with [NewSliceValueWriter].
type SliceValueWriter struct {
	ValueWriter
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// ShowAddr shows the address.
	// Default: true.
	ShowAddr bool
	// MaxLen is the maximum length of the slice.
	// Default: 0 (no limit).
	MaxLen int
}

// NewSliceValueWriter creates a new [SliceValueWriter] with default values.
func NewSliceValueWriter(vw ValueWriter) *SliceValueWriter {
	return &SliceValueWriter{
		ValueWriter: vw,
		ShowLen:     true,
		ShowCap:     true,
		ShowAddr:    true,
		MaxLen:      0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *SliceValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Slice {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	infos{
		showLen:  vw.ShowLen,
		len:      v.Len(),
		showCap:  vw.ShowCap,
		cap:      v.Cap(),
		showAddr: vw.ShowAddr,
		addr:     uint64(uintptr(v.UnsafePointer())),
	}.writeWithTrailingSpace(w)
	writeArray(c, w, st, v, vw.MaxLen, vw.ValueWriter)
	return true
}

// StringValueWriter is a [ValueWriter] that handles string values.
//
// It should be created with [NewStringValueWriter].
type StringValueWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// Quote quotes the string.
	// Default: true.
	Quote bool
	// MaxLen is the maximum length of the string.
	// Default: 0 (no limit).
	MaxLen int
}

// NewStringValueWriter creates a new [StringValueWriter] with default values.
func NewStringValueWriter() *StringValueWriter {
	return &StringValueWriter{
		ShowLen: true,
		Quote:   true,
		MaxLen:  0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StringValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.String {
		return false
	}
	s := v.String()
	writeStringValue(w, s, vw.ShowLen, vw.Quote, vw.MaxLen)
	return true
}

func writeStringValue(w io.Writer, s string, showLen bool, quote bool, maxLen int) {
	infos{
		showLen: showLen,
		len:     len(s),
	}.writeWithTrailingSpace(w)
	truncated := false
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
		truncated = true
	}
	if quote {
		writeQuote(w, s)
	} else {
		writeString(w, s)
	}
	if truncated {
		writeString(w, " ")
		writeTruncated(w)
	}
}

// StructValueWriter is a [ValueWriter] that handles struct values.
//
// It should be created with [NewStructValueWriter].
type StructValueWriter struct {
	ValueWriter
	// Unexported shows unexported fields.
	// Default: true.
	Unexported bool
}

// NewStructValueWriter creates a new [StructValueWriter] with default values.
func NewStructValueWriter(vw ValueWriter) *StructValueWriter {
	return &StructValueWriter{
		ValueWriter: vw,
		Unexported:  true,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StructValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct {
		return false
	}
	writeString(w, "{")
	fields := getStructFields(v.Type())
	if len(fields) > 0 {
		writeString(w, "\n")
		st.RunIndent(func(st *State) {
			for i, field := range fields {
				if !vw.Unexported && !field.IsExported() {
					continue
				}
				c.WriteIndent(w, st)
				writeString(w, field.Name)
				writeString(w, ": ")
				mustHandle(vw.ValueWriter(c, w, st, v.Field(i)))
				writeString(w, ",\n")
			}
		})
		c.WriteIndent(w, st)
	}
	writeString(w, "}")
	return true
}

var (
	structFieldsCacheLock sync.Mutex
	structFieldsCache     = map[reflect.Type][]reflect.StructField{}
)

func getStructFields(typ reflect.Type) []reflect.StructField {
	structFieldsCacheLock.Lock()
	fields, ok := structFieldsCache[typ]
	if !ok {
		fields = make([]reflect.StructField, typ.NumField())
		for i := range typ.NumField() {
			fields[i] = typ.Field(i)
		}
		structFieldsCache[typ] = fields
	}
	structFieldsCacheLock.Unlock()
	return fields
}

// UnsafePointerValueWriter is a [ValueWriter] that handles unsafe pointer values.
//
// It should be created with [NewUnsafePointerValueWriter].
type UnsafePointerValueWriter struct{}

// NewUnsafePointerValueWriter creates a new [UnsafePointerValueWriter].
func NewUnsafePointerValueWriter() *UnsafePointerValueWriter {
	return &UnsafePointerValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *UnsafePointerValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() != reflect.UnsafePointer {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	writeUintptr(w, uint64(uintptr(v.UnsafePointer())))
	return true
}

// ValueWriters is a list of [ValueWriter].
//
// They are tried in order until one handles the value.
type ValueWriters []ValueWriter

// WriteValue implements [ValueWriter].
func (vws ValueWriters) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	for _, vw := range vws {
		ok := vw(c, w, st, v)
		if ok {
			return true
		}
	}
	return false
}

// FilterValueWriter is a [ValueWriter] that calls the [ValueWriter] if the filter returns true.
//
// It should be created with [NewFilterValueWriter].
type FilterValueWriter struct {
	ValueWriter
	Filter func(v reflect.Value) bool
}

// NewFilterValueWriter creates a new [FilterValueWriter].
func NewFilterValueWriter(vw ValueWriter, f func(v reflect.Value) bool) *FilterValueWriter {
	return &FilterValueWriter{
		ValueWriter: vw,
		Filter:      f,
	}
}

// WriteValue implements [ValueWriter].
func (vw *FilterValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if !vw.Filter(v) {
		return false
	}
	return vw.ValueWriter(c, w, st, v)
}

// DefaultCommonValueWriter is the default [CommonValueWriter] created with [NewCommonValueWriter].
var DefaultCommonValueWriter = NewCommonValueWriter()

// CommonValueWriter is a [ValueWriter] with common [ValueWriter].
//
// Any [ValueWriter] can be configured, but it's not allowed to change the pointer value.
// Any [ValueWriter] can be set to nil in order to disable it.
// It is not allowed to updated the wrapped [ValueWriter].
//
// It should be created with [NewCommonValueWriter].
type CommonValueWriter struct {
	PanicRecover     *PanicRecoverValueWriter
	UnwrapInterface  *UnwrapInterfaceValueWriter
	Recursion        *RecursionValueWriter
	MaxDepth         *MaxDepthValueWriter
	TypeAndValue     *TypeAndValueWriter
	Type             *TypeValueWriter
	Kind             *KindValueWriter
	ReflectValue     *ReflectValueWriter
	Error            *ErrorValueWriter
	BytesHexDump     *BytesHexDumpValueWriter
	BytesableHexDump *BytesableHexDumpValueWriter
	Stringer         *StringerValueWriter
	ValueWriters     ValueWriters
}

// NewCommonValueWriter creates a new [CommonValueWriter] initialized with default values.
func NewCommonValueWriter() *CommonValueWriter {
	vw := &CommonValueWriter{}
	vw.PanicRecover = NewPanicRecoverValueWriter(vw.postPanicRecover)
	vw.UnwrapInterface = NewUnwrapInterfaceValueWriter(vw.postUnwrapInterface)
	vw.Recursion = NewRecursionValueWriter(vw.postRecursion)
	vw.MaxDepth = NewMaxDepthValueWriter(vw.postMaxDepth)
	vw.TypeAndValue = NewTypeAndValueWriter(vw.writeType, vw.postTypeAndValue)
	vw.Type = NewTypeValueWriter()
	vw.Kind = NewKindValueWriter(vw.loopback)
	vw.ReflectValue = NewReflectValueWriter(vw.loopback)
	vw.Error = NewErrorValueWriter()
	vw.BytesHexDump = NewBytesHexDumpValueWriter()
	vw.BytesableHexDump = NewBytesableHexDumpValueWriter()
	vw.Stringer = NewStringerValueWriter()
	vw.ValueWriters = ValueWriters{
		vw.reflectValue,
		vw.error,
		vw.bytesHexDump,
		vw.bytesableHexDump,
		vw.stringer,
	}
	return vw
}

// SetShowLen sets ShowLen on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowLen(show bool) {
	vw.Kind.BaseChan.ShowLen = show
	vw.Kind.BaseMap.ShowLen = show
	vw.Kind.BaseSlice.ShowLen = show
	vw.Kind.BaseString.ShowLen = show
	vw.BytesHexDump.ShowLen = show
	vw.BytesableHexDump.ShowLen = show
	vw.Stringer.ShowLen = show
}

// SetShowCap sets ShowCap on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowCap(show bool) {
	vw.Kind.BaseChan.ShowCap = show
	vw.Kind.BaseSlice.ShowCap = show
	vw.BytesHexDump.ShowCap = show
	vw.BytesableHexDump.ShowCap = show
}

// SetShowAddr sets ShowAddr on all [ValueWriter] that supports it.
func (vw *CommonValueWriter) SetShowAddr(show bool) {
	vw.Kind.BaseChan.ShowAddr = show
	vw.Kind.BaseFunc.ShowAddr = show
	vw.Kind.BaseMap.ShowAddr = show
	vw.Kind.BasePointer.ShowAddr = show
	vw.Kind.BaseSlice.ShowAddr = show
}

// ConfigureWithPrinter configures the [CommonValueWriter] with a [Printer].
//
// It sets the [MapValueWriter.SortKeysCmpDefault] to use [Printer.compare].
func (vw *CommonValueWriter) ConfigureWithPrinter(p *Printer) {
	vw.Kind.BaseMap.SortKeysCmpDefault = p.compare
}

// ConfigureTest configures the [CommonValueWriter] for testing.
//
// It makes the result deterministic.
// It sorts the keys of maps and disables the address.
func (vw *CommonValueWriter) ConfigureTest() {
	vw.Kind.BaseMap.SortKeys = true
	vw.SetShowAddr(false)
}

// WriteValue implements [ValueWriter].
func (vw *CommonValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.panicRecover(c, w, st, v)
}

func (vw *CommonValueWriter) panicRecover(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.PanicRecover == nil {
		return vw.postPanicRecover(c, w, st, v)
	}
	return vw.PanicRecover.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) postPanicRecover(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.unwrapInterface(c, w, st, v)
}

func (vw *CommonValueWriter) loopback(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.unwrapInterface(c, w, st, v)
}

func (vw *CommonValueWriter) unwrapInterface(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.UnwrapInterface == nil {
		return vw.postUnwrapInterface(c, w, st, v)
	}
	return vw.UnwrapInterface.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) postUnwrapInterface(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.recursion(c, w, st, v)
}

func (vw *CommonValueWriter) recursion(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.Recursion == nil {
		return vw.postRecursion(c, w, st, v)
	}
	return vw.Recursion.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) postRecursion(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.maxDepth(c, w, st, v)
}

func (vw *CommonValueWriter) maxDepth(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.MaxDepth == nil {
		return vw.postMaxDepth(c, w, st, v)
	}
	return vw.MaxDepth.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) postMaxDepth(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.typeAndValue(c, w, st, v)
}

func (vw *CommonValueWriter) writeType(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.Type.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) typeAndValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.TypeAndValue == nil || vw.Type == nil {
		return vw.postTypeAndValue(c, w, st, v)
	}
	return vw.TypeAndValue.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) postTypeAndValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	return vw.internal(c, w, st, v)
}

func (vw *CommonValueWriter) internal(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.ValueWriters.WriteValue(c, w, st, v) {
		return true
	}
	return vw.kind(c, w, st, v)
}

func (vw *CommonValueWriter) kind(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.Kind == nil {
		return false
	}
	return vw.Kind.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) reflectValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.ReflectValue == nil {
		return false
	}
	return vw.ReflectValue.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) error(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.Error == nil {
		return false
	}
	return vw.Error.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) bytesHexDump(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.BytesHexDump == nil {
		return false
	}
	return vw.BytesHexDump.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) bytesableHexDump(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.BytesableHexDump == nil {
		return false
	}
	return vw.BytesableHexDump.WriteValue(c, w, st, v)
}

func (vw *CommonValueWriter) stringer(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.Stringer == nil {
		return false
	}
	return vw.Stringer.WriteValue(c, w, st, v)
}

// PanicRecoverValueWriter is a [ValueWriter] that recovers from panics.
//
// It should be created with [NewPanicRecoverValueWriter].
type PanicRecoverValueWriter struct {
	ValueWriter
}

// NewPanicRecoverValueWriter creates a new [PanicRecoverValueWriter].
func NewPanicRecoverValueWriter(vw ValueWriter) *PanicRecoverValueWriter {
	return &PanicRecoverValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *PanicRecoverValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) (handled bool) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		handled = true
		_, _ = writeStringErr(w, "<panic>: ")
		switch r := r.(type) {
		case string:
			_, _ = writeStringErr(w, r)
		case error:
			_, _ = writeStringErr(w, r.Error())
		default:
			_, _ = fmt.Fprint(w, r)
		}
		_, _ = writeStringErr(w, "\n")
		// TODO print stack trace
	}()
	return vw.ValueWriter(c, w, st, v)
}

// UnwrapInterfaceValueWriter is a [ValueWriter] that unwraps interface values.
//
// It should be created with [NewUnwrapInterfaceValueWriter].
type UnwrapInterfaceValueWriter struct {
	ValueWriter
}

// NewUnwrapInterfaceValueWriter creates a new [UnwrapInterfaceValueWriter].
func NewUnwrapInterfaceValueWriter(vw ValueWriter) *UnwrapInterfaceValueWriter {
	return &UnwrapInterfaceValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *UnwrapInterfaceValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Kind() == reflect.Interface {
		if checkNil(w, v) {
			return true
		}
		v = v.Elem()
	}
	return vw.ValueWriter(c, w, st, v)
}

// RecursionValueWriter is a [ValueWriter] that prevents recursion.
//
// It should be created with [NewRecursionValueWriter].
type RecursionValueWriter struct {
	ValueWriter
}

// NewRecursionValueWriter creates a new [RecursionValueWriter].
func NewRecursionValueWriter(vw ValueWriter) *RecursionValueWriter {
	return &RecursionValueWriter{
		ValueWriter: vw,
	}
}

// WriteValue implements [ValueWriter].
func (vw *RecursionValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Only handles pointer kinds.
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
	default:
		return vw.ValueWriter(c, w, st, v)
	}
	vp := v.Pointer()
	if slices.Contains(st.Visited, vp) {
		writeString(w, "<recursion>")
		return true
	}
	var handled bool
	st.RunVisited(vp, func(st *State) {
		handled = vw.ValueWriter(c, w, st, v)
	})
	return handled
}

// MaxDepthValueWriter is a [ValueWriter] that limits the depth.
//
// It should be created with [NewMaxDepthValueWriter].
type MaxDepthValueWriter struct {
	ValueWriter
	// Max is the maximum depth.
	// Default: 0 (no limit).
	Max int
}

// NewMaxDepthValueWriter creates a new [MaxDepthValueWriter].
func NewMaxDepthValueWriter(vw ValueWriter) *MaxDepthValueWriter {
	return &MaxDepthValueWriter{
		ValueWriter: vw,
		Max:         0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MaxDepthValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if vw.Max <= 0 {
		return vw.ValueWriter(c, w, st, v)
	}
	if st.Depth >= vw.Max {
		writeString(w, "<max depth>")
		return true
	}
	var handled bool
	st.RunDepth(func(st *State) {
		handled = vw.ValueWriter(c, w, st, v)
	})
	return handled
}

// TypeAndValueWriter is a [ValueWriter] that writes the type and the value.
//
// It should be created with [NewTypeAndValueWriter].
type TypeAndValueWriter struct {
	Type  ValueWriter
	Value ValueWriter
}

// NewTypeAndValueWriter creates a new [TypeAndValueWriter].
func NewTypeAndValueWriter(t, v ValueWriter) *TypeAndValueWriter {
	return &TypeAndValueWriter{
		Type:  t,
		Value: v,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeAndValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	writeString(w, "(")
	mustHandle(vw.Type(c, w, st, v))
	writeString(w, ") ")
	mustHandle(vw.Value(c, w, st, v))
	return true
}

// TypeValueWriter is a [ValueWriter] that writes the type.
//
// It should be created with [NewTypeValueWriter].
type TypeValueWriter struct {
	// Stringer converts the [reflect.Type] to a string.
	// Default: [reflect.Type.String].
	Stringer func(reflect.Type) string
}

// NewTypeValueWriter creates a new [TypeValueWriter] with default values.
func NewTypeValueWriter() *TypeValueWriter {
	return &TypeValueWriter{
		Stringer: reflect.Type.String,
	}
}

// WriteValue implements [ValueWriter].
func (vw *TypeValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	writeString(w, vw.Stringer(v.Type()))
	return true
}

var typeReflectValue = reflect.TypeFor[reflect.Value]()

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
func (vw *ReflectValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Type() != typeReflectValue {
		return false
	}
	if !v.CanInterface() {
		writeString(w, "<unexported>")
		return true
	}
	rv := v.Interface().(reflect.Value) //nolint:forcetypeassert // Checked above.
	writeArrow(w)
	if checkInvalid(w, rv) {
		return true
	}
	mustHandle(vw.ValueWriter(c, w, st, rv))
	return true
}

var typeError = reflect.TypeFor[error]()

// ErrorValueWriter is a [ValueWriter] that handles errors and write them with error.Error.
//
// It should be created with [NewErrorValueWriter].
type ErrorValueWriter struct{}

// NewErrorValueWriter creates a new [ErrorValueWriter].
func NewErrorValueWriter() *ErrorValueWriter {
	return &ErrorValueWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *ErrorValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if !v.Type().Implements(typeError) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	writeArrowWrappedString(w, ".Error() ")
	writeQuote(w, err.Error())
	return true
}

var bytesType = reflect.TypeFor[[]byte]()

// BytesHexDumpValueWriter is a [ValueWriter] that handles []byte and writes them with [hex.Dumper].
//
// It should be created with [NewBytesHexDumpValueWriter].
type BytesHexDumpValueWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// MaxLen is the maximum length of the bytes.
	// Default: 0 (no limit).
	MaxLen int
}

// NewBytesHexDumpValueWriter creates a new [BytesHexDumpValueWriter].
func NewBytesHexDumpValueWriter() *BytesHexDumpValueWriter {
	return &BytesHexDumpValueWriter{
		ShowLen: true,
		ShowCap: true,
		MaxLen:  0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *BytesHexDumpValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if v.Type() != bytesType {
		return false
	}
	if checkNil(w, v) {
		return true
	}
	b := v.Bytes()
	writeBytesHexDumpCommon(c, w, st, b, vw.ShowLen, vw.ShowCap, vw.MaxLen)
	return true
}

// Bytesable is an interface that can return a []byte.
type Bytesable interface {
	Bytes() []byte
}

var bytesableType = reflect.TypeFor[Bytesable]()

// BytesableHexDumpValueWriter is a [ValueWriter] that handles [Bytesable] and writes thems with [hex.Dumper].
//
// It should be created with [NewBytesableHexDumpValueWriter].
type BytesableHexDumpValueWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowCap shows the cap.
	// Default: true.
	ShowCap bool
	// MaxLen is the maximum length of the bytes.
	// Default: 0 (no limit).
	MaxLen int
}

// NewBytesableHexDumpValueWriter creates a new [BytesableHexDumpValueWriter].
func NewBytesableHexDumpValueWriter() *BytesableHexDumpValueWriter {
	return &BytesableHexDumpValueWriter{
		ShowLen: true,
		ShowCap: true,
		MaxLen:  0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *BytesableHexDumpValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if !v.Type().Implements(bytesableType) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if v.Type() == typeReflectValue {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	br := v.Interface().(Bytesable) //nolint:forcetypeassert // Checked above.
	b := br.Bytes()
	writeArrowWrappedString(w, ".Bytes() ")
	if b == nil {
		writeNil(w)
		return true
	}
	writeBytesHexDumpCommon(c, w, st, b, vw.ShowLen, vw.ShowCap, vw.MaxLen)
	return true
}

func writeBytesHexDumpCommon(c *Config, w io.Writer, st *State, b []byte, showLen bool, showCap bool, maxLen int) {
	infos{
		showLen: showLen,
		len:     len(b),
		showCap: showCap,
		cap:     cap(b),
	}.writeWithTrailingSpace(w)
	truncated := false
	if maxLen > 0 && len(b) > maxLen {
		b = b[:maxLen]
		truncated = true
	}
	writeString(w, "\n")
	st.RunIndent(func(st *State) {
		iw := GetIndentWriter(c, w, st, false)
		defer iw.Release()
		d := hex.Dumper(iw)
		mustWrite(d.Write(b))
		must(d.Close())
		if truncated {
			c.WriteIndent(w, st)
			writeTruncated(w)
		}
	})
}

var typeStringer = reflect.TypeFor[fmt.Stringer]()

// StringerValueWriter is a [ValueWriter] that handles [fmt.Stringer].
//
// It should be created with [NewStringerValueWriter].
type StringerValueWriter struct {
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// Quote quotes the string.
	// Default: true.
	Quote bool
	// MaxLen is the maximum length of the string.
	// Default: 0 (no limit).
	MaxLen int
}

// NewStringerValueWriter creates a new [StringerValueWriter].
func NewStringerValueWriter() *StringerValueWriter {
	return &StringerValueWriter{
		ShowLen: true,
		Quote:   true,
		MaxLen:  0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *StringerValueWriter) WriteValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if !v.Type().Implements(typeStringer) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	if v.Type() == typeReflectValue {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	sr := v.Interface().(fmt.Stringer) //nolint:forcetypeassert // Checked above.
	s := sr.String()
	writeArrowWrappedString(w, ".String() ")
	writeStringValue(w, s, vw.ShowLen, vw.Quote, vw.MaxLen)
	return true
}

type formatter struct {
	printer *Printer
	value   any
}

func (ft *formatter) Format(f fmt.State, verb rune) {
	ft.printer.Write(f, ft.value)
}

// IndentWriter is a [io.Writer] that indents.
//
// It is exposed for internal use.
// It should not be used outside of this package.
type IndentWriter struct {
	writer   io.Writer
	config   *Config
	state    *State
	indented bool
}

// NewIndentWriter creates a new [IndentWriter].
func NewIndentWriter(c *Config, w io.Writer, st *State, indented bool) *IndentWriter {
	iw := &IndentWriter{}
	iw.init(c, w, st, indented)
	return iw
}

func (iw *IndentWriter) init(c *Config, w io.Writer, st *State, indented bool) {
	iw.writer = w
	iw.config = c
	iw.state = st
	iw.indented = indented
}

func (iw *IndentWriter) reset() {
	iw.writer = nil
	iw.config = nil
	iw.state = nil
	iw.indented = false
}

// Write implements [io.Writer].
func (iw *IndentWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if !iw.indented {
			nn, err := writeIndentErr(iw.writer, iw.config.Indent, iw.state.Indent)
			n += nn
			if err != nil {
				return n, err
			}
			iw.indented = true
		}
		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			i = len(p)
		} else {
			i++
			iw.indented = false
		}
		nn, err := iw.writer.Write(p[:i])
		n += nn
		if err != nil {
			return n, err //nolint:wrapcheck // The error is not wrapped.
		}
		p = p[i:]
	}
	return n, nil
}

var indentWriterPool = &sync.Pool{
	New: func() any {
		return &IndentWriter{}
	},
}

// GetIndentWriter returns a [IndentWriter] from a pool.
//
// The caller must call [IndentWriter.Release] after using it.
func GetIndentWriter(c *Config, w io.Writer, st *State, indented bool) *IndentWriter {
	iw := indentWriterPool.Get().(*IndentWriter) //nolint:forcetypeassert // The pool only contains *indentWriter.
	iw.init(c, w, st, indented)
	return iw
}

// Release releases the [IndentWriter] to the pool.
func (iw *IndentWriter) Release() {
	iw.reset()
	indentWriterPool.Put(iw)
}

var (
	indentCacheLock sync.Mutex
	indentCache     = map[string][]byte{}
)

func getIndent(s string, n int) []byte {
	indentCacheLock.Lock()
	defer indentCacheLock.Unlock()
	b := indentCache[s]
	l := len(s) * n
	if len(b) < l {
		b = bytes.Repeat([]byte(s), n)
		indentCache[s] = b
	}
	return b[:l]
}

// WriteIndent writes the indent s (n times) to the [io.Writer].
func WriteIndent(w io.Writer, s string, n int) {
	mustWrite(writeIndentErr(w, s, n))
}

func writeIndentErr(w io.Writer, s string, n int) (int, error) {
	if n <= 0 {
		return 0, nil
	}
	if n == 1 {
		return writeStringErr(w, s)
	}
	return w.Write(getIndent(s, n)) //nolint:wrapcheck // The error is not wrapped.
}

func checkNil(w io.Writer, v reflect.Value) bool {
	if v.IsNil() {
		writeNil(w)
		return true
	}
	return false
}

func writeNil(w io.Writer) {
	writeString(w, "<nil>")
}

func writeArrow(w io.Writer) {
	writeString(w, "=> ")
}

func writeArrowWrappedString(w io.Writer, s string) {
	writeArrow(w)
	writeString(w, s)
	writeArrow(w)
}

func writeTruncated(w io.Writer) {
	writeString(w, "<truncated>")
}

type infos struct {
	showLen  bool
	len      int
	showCap  bool
	cap      int
	showAddr bool
	addr     uint64
}

func (i infos) write(w io.Writer) bool {
	if !i.showLen && !i.showCap && !i.showAddr {
		return false
	}
	writeString(w, "(")
	wrote := false
	if i.showLen {
		writeString(w, "len=")
		mustWrite(strconvio.WriteInt(w, int64(i.len), 10))
		wrote = true
	}
	if i.showCap {
		if wrote {
			writeString(w, " ")
		}
		writeString(w, "cap=")
		mustWrite(strconvio.WriteInt(w, int64(i.cap), 10))
		wrote = true
	}
	if i.showAddr {
		if wrote {
			writeString(w, " ")
		}
		writeString(w, "addr=")
		writeUintptr(w, i.addr)
	}
	writeString(w, ")")
	return true
}

func (i infos) writeWithTrailingSpace(w io.Writer) {
	if i.write(w) {
		writeString(w, " ")
	}
}

func writeString(w io.Writer, s string) {
	mustWrite(writeStringErr(w, s))
}

func writeQuote(w io.Writer, s string) {
	mustWrite(strconvio.WriteQuote(w, s))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustWrite(_ int, err error) {
	must(err)
}

func mustHandle(h bool) {
	if !h {
		panic("not handled")
	}
}
