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
	"sync"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/strconvio"
)

// Write writes the value to the writer with [DefaultConfig].
func Write(w io.Writer, vi any) {
	DefaultConfig.Write(w, vi)
}

// String returns the value as a string with [DefaultConfig].
func String(vi any) string {
	return DefaultConfig.String(vi)
}

// Formatter returns a [fmt.Formatter] for the value with [DefaultConfig].
func Formatter(vi any) fmt.Formatter {
	return DefaultConfig.Formatter(vi)
}

// DefaultConfig is the default [Config].
var DefaultConfig = NewConfig()

// Config is a configuration used to pretty print values.
//
// It should be created with [NewConfig].
type Config struct {
	// PanicRecover recovers from panics and writes it to the writer.
	// Default: true.
	PanicRecover bool
	// TypeFullName prints the full type name.
	// Default: false.
	TypeFullName bool
	// Indent is the string used to indent.
	// Default: "\t".
	Indent string
	// MaxDepth is the maximum depth.
	// Default: 0 (no limit).
	MaxDepth int
	// StringMaxLen is the maximum length of strings.
	// Default: 0 (no limit).
	StringMaxLen int
	// SliceMaxLen is the maximum length of slices and arrays.
	// Default: 0 (no limit).
	SliceMaxLen int
	// MapSortKeys sorts map keys.
	// Default: false.
	MapSortKeys bool
	// MapMaxLen is the maximum length of maps.
	// Default: 0 (no limit).
	MapMaxLen int
	// StructUnexported prints unexported fields of structs.
	// Default: true.
	StructUnexported bool
	// ValueWriters is the list of ValueWriter used to write values.
	// Default: reflect.Value, error, []byte, interface{ Bytes() []byte }, fmt.Stringer.
	ValueWriters []ValueWriter
}

// NewConfig creates a new [Config] initialized with default values.
func NewConfig() *Config {
	return &Config{
		PanicRecover:     true,
		Indent:           "\t",
		StructUnexported: true,
		ValueWriters: []ValueWriter{
			NewReflectValueValueWriter(),
			NewErrorValueWriter(),
			NewBytesHexValueWriter(0),
			NewByteserHexValueWriter(0),
			NewStringerValueWriter(0),
		},
	}
}

// Write writes the value to the writer.
func (c *Config) Write(w io.Writer, vi any) {
	v := reflect.ValueOf(vi)
	c.write(w, v)
}

func (c *Config) write(w io.Writer, v reflect.Value) {
	stItf := statePool.Get()
	defer statePool.Put(stItf)
	st := stItf.(*State) //nolint:forcetypeassert // The pool only contains *State.
	st.reset()
	c.WriteTypeAndValue(w, st, v)
}

var bufPool = &bufpool.Pool{}

// String returns the value as a string.
func (c *Config) String(vi any) string {
	v := reflect.ValueOf(vi)
	return c.string(v)
}

func (c *Config) string(v reflect.Value) string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	c.write(buf, v)
	return buf.String()
}

// Formatter returns a [fmt.Formatter] for the value.
func (c *Config) Formatter(vi any) fmt.Formatter {
	return &formatter{
		config: c,
		value:  vi,
	}
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

// WriteIndent writes the indentation to the writer.
func (c *Config) WriteIndent(w io.Writer, st *State) {
	if st.Indent <= 0 {
		return
	}
	if st.Indent == 1 {
		_, _ = WriteString(w, c.Indent)
		return
	}
	_, _ = w.Write(getIndent(c.Indent, st.Indent))
}

func (c *Config) runRecursion(w io.Writer, st *State, v reflect.Value, f func(st *State)) {
	vp := v.Pointer()
	if slices.Contains(st.Visited, vp) {
		_, _ = WriteString(w, "<recursion>")
		return
	}
	st.RunVisited(vp, f)
}

// WriteTypeAndValue writes the type and value to the writer.
//
// It writes "(TYPE) VALUE".
func (c *Config) WriteTypeAndValue(w io.Writer, st *State, v reflect.Value) {
	if c.PanicRecover {
		defer func() {
			c.checkRecover(w, recover())
		}()
	}
	if !v.IsValid() {
		WriteNil(w)
		return
	}
	if v.Kind() == reflect.Interface {
		c.WriteTypeAndValue(w, st, v.Elem())
		return
	}
	if c.MaxDepth > 0 && st.Depth >= c.MaxDepth {
		_, _ = WriteString(w, "<max depth>")
		return
	}
	st.RunDepth(func(st *State) {
		_, _ = WriteString(w, "(")
		c.WriteType(w, v.Type())
		_, _ = WriteString(w, ") ")
		c.WriteValue(w, st, v)
	})
}

func (c *Config) checkRecover(w io.Writer, r any) {
	if r == nil {
		return
	}
	_, _ = WriteString(w, "<panic>: ")
	switch r := r.(type) {
	case string:
		_, _ = WriteString(w, r)
	case error:
		_, _ = WriteString(w, r.Error())
	default:
		_, _ = fmt.Fprint(w, r)
	}
	_, _ = WriteString(w, "\n")
}

// WriteType writes the type to the writer.
func (c *Config) WriteType(w io.Writer, typ reflect.Type) {
	var s string
	if c.TypeFullName {
		s = reflectutil.TypeFullName(typ)
	} else {
		s = typ.String()
	}
	_, _ = WriteString(w, s)
}

// WriteValue writes the value to the writer.
//
// It checks if any of the [ValueWriter] can handle the value, then call [Config.WriteValueDefault].
func (c *Config) WriteValue(w io.Writer, st *State, v reflect.Value) {
	if c.writeValueWithValueWriter(w, st, v) {
		return
	}
	c.WriteValueDefault(w, st, v)
}

func (c *Config) writeValueWithValueWriter(w io.Writer, st *State, v reflect.Value) bool {
	for _, vw := range c.ValueWriters {
		ok := vw(c, w, st, v)
		if ok {
			return true
		}
	}
	return false
}

// WriteValueDefault writes the value to the writer with the default behavior.
//
// It skips all the [ValueWriter].
//
//nolint:gocyclo // We need to handle all kinds.
func (c *Config) WriteValueDefault(w io.Writer, st *State, v reflect.Value) {
	switch v.Kind() { //nolint:exhaustive // All kinds are handled, Invalid and Interface should not happen.
	case reflect.Bool:
		c.writeBool(w, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		c.writeInt(w, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		c.writeUint(w, v)
	case reflect.Uintptr:
		c.writeUintptr(w, v)
	case reflect.Float32, reflect.Float64:
		c.writeFloat(w, v)
	case reflect.Complex64, reflect.Complex128:
		c.writeComplex(w, v)
	case reflect.String:
		c.writeString(w, v)
	case reflect.Chan:
		c.writeChan(w, v)
	case reflect.Func:
		c.writeFunc(w, v)
	case reflect.Pointer:
		c.writePointer(w, st, v)
	case reflect.UnsafePointer:
		c.writeUnsafePointer(w, v)
	case reflect.Array:
		c.writeArray(w, st, v)
	case reflect.Slice:
		c.writeSlice(w, st, v)
	case reflect.Map:
		c.writeMap(w, st, v)
	case reflect.Struct:
		c.writeStruct(w, st, v)
	default:
		panic(fmt.Sprintf("unexpected kind %s", v.Kind()))
	}
}

func (c *Config) writeBool(w io.Writer, v reflect.Value) {
	_, _ = strconvio.WriteBool(w, v.Bool())
}

func (c *Config) writeInt(w io.Writer, v reflect.Value) {
	_, _ = strconvio.WriteInt(w, v.Int(), 10)
}

func (c *Config) writeUint(w io.Writer, v reflect.Value) {
	_, _ = strconvio.WriteUint(w, v.Uint(), 10)
}

func (c *Config) writeUintptr(w io.Writer, v reflect.Value) {
	_, _ = WriteString(w, "0x")
	_, _ = strconvio.WriteUint(w, v.Uint(), 16)
}

func (c *Config) writeFloat(w io.Writer, v reflect.Value) {
	bitSize := v.Type().Bits()
	_, _ = strconvio.WriteFloat(w, v.Float(), 'g', -1, bitSize)
}

func (c *Config) writeComplex(w io.Writer, v reflect.Value) {
	_, _ = fmt.Fprintf(w, "%g", v.Complex())
}

func (c *Config) writeString(w io.Writer, v reflect.Value) {
	s := v.String()
	writeStringValue(w, s, c.StringMaxLen)
}

func writeStringValue(w io.Writer, s string, maxLen int) {
	_, _ = WriteString(w, "(len=")
	_, _ = strconvio.WriteInt(w, int64(len(s)), 10)
	_, _ = WriteString(w, ") ")
	truncated := false
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
		truncated = true
	}
	_, _ = strconvio.WriteQuote(w, s)
	if truncated {
		_, _ = WriteString(w, " ")
		writeTruncated(w)
	}
}

func (c *Config) writeChan(w io.Writer, v reflect.Value) {
	if v.IsNil() {
		WriteNil(w)
		return
	}
	writeLenCapReflect(w, v)
}

func (c *Config) writeFunc(w io.Writer, v reflect.Value) {
	if v.IsNil() {
		WriteNil(w)
		return
	}
	name := runtime.FuncForPC(v.Pointer()).Name()
	_, _ = WriteString(w, name)
}

func (c *Config) writePointer(w io.Writer, st *State, v reflect.Value) {
	c.runRecursion(w, st, v, func(st *State) {
		WriteArrow(w)
		c.WriteTypeAndValue(w, st, v.Elem())
	})
}

func (c *Config) writeUnsafePointer(w io.Writer, v reflect.Value) {
	_, _ = WriteString(w, "0x")
	_, _ = strconvio.WriteUint(w, uint64(uintptr(v.UnsafePointer())), 16)
}

func (c *Config) writeArray(w io.Writer, st *State, v reflect.Value) {
	l := v.Len()
	truncated := false
	if c.SliceMaxLen > 0 && l > c.SliceMaxLen {
		l = c.SliceMaxLen
		truncated = true
	}
	_, _ = WriteString(w, "{\n")
	if v.Len() > 0 {
		st.RunIndent(func(st *State) {
			for i := range l {
				c.WriteIndent(w, st)
				c.WriteTypeAndValue(w, st, v.Index(i))
				_, _ = WriteString(w, ",\n")
			}
			if truncated {
				c.WriteIndent(w, st)
				writeTruncated(w)
				_, _ = WriteString(w, "\n")
			}
		})
	}
	c.WriteIndent(w, st)
	_, _ = WriteString(w, "}")
}

func (c *Config) writeSlice(w io.Writer, st *State, v reflect.Value) {
	if v.IsNil() {
		WriteNil(w)
		return
	}
	c.runRecursion(w, st, v, func(st *State) {
		writeLenCapReflect(w, v)
		_, _ = WriteString(w, " ")
		c.writeArray(w, st, v)
	})
}

func (c *Config) writeMap(w io.Writer, st *State, v reflect.Value) {
	if v.IsNil() {
		WriteNil(w)
		return
	}
	c.runRecursion(w, st, v, func(st *State) {
		_, _ = WriteString(w, "(len=")
		_, _ = strconvio.WriteInt(w, int64(v.Len()), 10)
		_, _ = WriteString(w, ") {\n")
		if v.Len() > 0 {
			st.RunIndent(func(st *State) {
				if c.MapSortKeys {
					c.writeMapSorted(w, st, v)
				} else {
					c.writeMapUnsorted(w, st, v)
				}
			})
		}
		c.WriteIndent(w, st)
		_, _ = WriteString(w, "}")
	})
}

func (c *Config) writeMapSorted(w io.Writer, st *State, v reflect.Value) {
	keys := c.getSortedMapKeys(v)
	for i, key := range keys {
		ok := c.writeMapEntry(w, st, key, v.MapIndex(key), i)
		if !ok {
			break
		}
	}
}

func (c *Config) getSortedMapKeys(v reflect.Value) []reflect.Value {
	keys := v.MapKeys()
	c.sortMapKeys(v.Type().Key(), keys)
	return keys
}

func (c *Config) sortMapKeys(typ reflect.Type, vs []reflect.Value) {
	cmpFunc := c.getMapKeysSortCmp(typ)
	slices.SortFunc(vs, cmpFunc)
}

func (c *Config) getMapKeysSortCmp(typ reflect.Type) func(a, b reflect.Value) int {
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
	default:
		return func(a, b reflect.Value) int {
			return cmp.Compare(c.string(a), c.string(b))
		}
	}
}

func (c *Config) writeMapUnsorted(w io.Writer, st *State, v reflect.Value) {
	if v.CanInterface() {
		c.writeMapUnsortedExported(w, st, v)
	} else {
		c.writeMapUnsortedUnexported(w, st, v)
	}
}

var typeInterface = reflect.TypeFor[any]()

var reflectValuePool = &sync.Pool{
	New: func() any {
		return reflect.New(typeInterface).Elem()
	},
}

func (c *Config) writeMapUnsortedExported(w io.Writer, st *State, v reflect.Value) {
	iter := v.MapRange()
	keyItf := reflectValuePool.Get()
	valueItf := reflectValuePool.Get()
	key := keyItf.(reflect.Value)     //nolint:forcetypeassert // The pool only contains *State.
	value := valueItf.(reflect.Value) //nolint:forcetypeassert // The pool only contains *State.
	defer func() {
		key.SetZero()
		value.SetZero()
		reflectValuePool.Put(keyItf)
		reflectValuePool.Put(valueItf)
	}()
	for i := 0; iter.Next(); i++ {
		key.SetIterKey(iter)
		value.SetIterValue(iter)
		ok := c.writeMapEntry(w, st, key, value, i)
		if !ok {
			break
		}
	}
}

func (c *Config) writeMapUnsortedUnexported(w io.Writer, st *State, v reflect.Value) {
	iter := v.MapRange()
	for i := 0; iter.Next(); i++ {
		key := iter.Key()
		value := iter.Value()
		ok := c.writeMapEntry(w, st, key, value, i)
		if !ok {
			break
		}
	}
}

func (c *Config) writeMapEntry(w io.Writer, st *State, key reflect.Value, value reflect.Value, i int) bool {
	c.WriteIndent(w, st)
	if c.MapMaxLen > 0 && i >= c.MapMaxLen {
		writeTruncated(w)
		_, _ = WriteString(w, "\n")
		return false
	}
	c.WriteTypeAndValue(w, st, key)
	_, _ = WriteString(w, ": ")
	c.WriteTypeAndValue(w, st, value)
	_, _ = WriteString(w, ",\n")
	return true
}

func (c *Config) writeStruct(w io.Writer, st *State, v reflect.Value) {
	_, _ = WriteString(w, "{\n")
	st.RunIndent(func(st *State) {
		fields := getStructFields(v.Type())
		for i, field := range fields {
			if !c.StructUnexported && !field.IsExported() {
				continue
			}
			c.WriteIndent(w, st)
			_, _ = WriteString(w, field.Name)
			_, _ = WriteString(w, ": ")
			c.WriteTypeAndValue(w, st, v.Field(i))
			_, _ = WriteString(w, ",\n")
		}
	})
	c.WriteIndent(w, st)
	_, _ = WriteString(w, "}")
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

var statePool = &sync.Pool{
	New: func() any {
		return &State{}
	},
}

// State represents the state of the pretty printer.
//
// Functions must restore the original state when they return.
type State struct {
	Depth   int
	Indent  int
	Visited []uintptr
}

// RunDepth runs the function with increased depth and restores the original depth after.
func (st *State) RunDepth(f func(st *State)) {
	st.Depth++
	defer func() {
		st.Depth--
	}()
	f(st)
}

// RunIndent runs the function with increased indentation and restores the original indentation after.
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

// ValueWriter is a function that writes a value.
// It can be used to override the default behavior.
//
// It returns true if it handled the value, false otherwise.
//
// Implementations must check [reflect.Value.CanInterface] before using [reflect.Value.Interface].
type ValueWriter func(c *Config, w io.Writer, st *State, v reflect.Value) bool

var typeReflectValue = reflect.TypeFor[reflect.Value]()

// NewReflectValueValueWriter returns a [ValueWriter] that writes [reflect.Value].
func NewReflectValueValueWriter() ValueWriter {
	return writeReflectValue
}

func writeReflectValue(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	if v.Type() != typeReflectValue {
		return false
	}
	rv := v.Interface().(reflect.Value) //nolint:forcetypeassert // Checked above.
	WriteArrow(w)
	c.WriteTypeAndValue(w, st, rv)
	return true
}

var typeError = reflect.TypeFor[error]()

// NewErrorValueWriter returns a [ValueWriter] that writes error.
func NewErrorValueWriter() ValueWriter {
	return writeError
}

func writeError(c *Config, w io.Writer, st *State, v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	if !v.Type().Implements(typeError) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	writeArrowWrappedString(w, ".Error() ")
	_, _ = strconvio.WriteQuote(w, err.Error())
	return true
}

var bytesType = reflect.TypeFor[[]byte]()

// NewBytesHexValueWriter returns a [ValueWriter] that writes []byte with [hex.Dumper].
func NewBytesHexValueWriter(maxLen int) ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		return writeBytesHex(c, w, st, v, maxLen)
	}
}

func writeBytesHex(c *Config, w io.Writer, st *State, v reflect.Value, maxLen int) bool {
	if v.Type() != bytesType {
		return false
	}
	if v.IsNil() {
		WriteNil(w)
		return true
	}
	writeLenCapReflect(w, v)
	b := v.Bytes()
	writeBytesHexCommon(c, w, st, b, maxLen)
	return true
}

type byteser interface {
	Bytes() []byte
}

var byteserType = reflect.TypeFor[byteser]()

// NewByteserHexValueWriter returns a [ValueWriter] that writes interface { Bytes() []byte } with [hex.Dumper].
func NewByteserHexValueWriter(maxLen int) ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		return writeByteserHex(c, w, st, v, maxLen)
	}
}

func writeByteserHex(c *Config, w io.Writer, st *State, v reflect.Value, maxLen int) bool {
	if !v.CanInterface() {
		return false
	}
	if !v.Type().Implements(byteserType) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	br := v.Interface().(byteser) //nolint:forcetypeassert // Checked above.
	b := br.Bytes()
	writeArrowWrappedString(w, ".Bytes() ")
	if b == nil {
		WriteNil(w)
		return true
	}
	writeLenCap(w, len(b), cap(b))
	writeBytesHexCommon(c, w, st, b, maxLen)
	return true
}

func writeBytesHexCommon(c *Config, w io.Writer, st *State, b []byte, maxLen int) {
	truncated := false
	if maxLen > 0 && len(b) > maxLen {
		b = b[:maxLen]
		truncated = true
	}
	_, _ = WriteString(w, "\n")
	st.RunIndent(func(st *State) {
		iw := GetIndentWriter(w, c, st, false)
		defer iw.Release()
		d := hex.Dumper(iw)
		_, _ = d.Write(b)
		_ = d.Close()
		if truncated {
			c.WriteIndent(w, st)
			writeTruncated(w)
		}
	})
}

var typeStringer = reflect.TypeFor[fmt.Stringer]()

// NewStringerValueWriter returns a [ValueWriter] that writes [fmt.Stringer].
func NewStringerValueWriter(maxLen int) ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		return writeStringer(w, v, maxLen)
	}
}

func writeStringer(w io.Writer, v reflect.Value, maxLen int) bool {
	if !v.CanInterface() {
		return false
	}
	if !v.Type().Implements(typeStringer) {
		return false
	}
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return false
	}
	sr := v.Interface().(fmt.Stringer) //nolint:forcetypeassert // Checked above.
	s := sr.String()
	writeArrowWrappedString(w, ".String() ")
	writeStringValue(w, s, maxLen)
	return true
}

// NewFilterValueWriter returns a [ValueWriter] that calls the provided [ValueWriter] if f returns true.
//
// It allows to enable/disable a [ValueWriter] for specific values/types.
func NewFilterValueWriter(vw ValueWriter, f func(v reflect.Value) bool) ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		return writeFilter(c, w, st, v, vw, f)
	}
}

func writeFilter(c *Config, w io.Writer, st *State, v reflect.Value, vw ValueWriter, f func(v reflect.Value) bool) bool {
	if !f(v) {
		return false
	}
	return vw(c, w, st, v)
}

type formatter struct {
	config *Config
	value  any
}

func (ft *formatter) Format(f fmt.State, verb rune) {
	ft.config.Write(f, ft.value)
}

// indentWriter is a io.Writer that indents.
type indentWriter struct {
	writer   io.Writer
	config   *Config
	state    *State
	indented bool
}

func (iw *indentWriter) Write(p []byte) (int, error) {
	l := len(p)
	for len(p) > 0 {
		if !iw.indented {
			iw.config.WriteIndent(iw.writer, iw.state)
			iw.indented = true
		}
		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			i = len(p)
		} else {
			i++
			iw.indented = false
		}
		_, err := iw.writer.Write(p[:i])
		if err != nil {
			return 0, err //nolint:wrapcheck // The error is not wrapped.
		}
		p = p[i:]
	}
	return l, nil
}

var indentWriterPool = &sync.Pool{
	New: func() any {
		return &indentWriter{}
	},
}

// GetIndentWriter returns an indentWriter from the pool.
//
// It is exposed for internal use.
// It should not be used outside of this package.
func GetIndentWriter(w io.Writer, c *Config, st *State, indented bool) *indentWriter {
	iw := indentWriterPool.Get().(*indentWriter) //nolint:forcetypeassert // The pool only contains *indentWriter.
	iw.writer = w
	iw.config = c
	iw.state = st
	iw.indented = indented
	return iw
}

// Release releases the indentWriter to the pool.
func (iw *indentWriter) Release() {
	iw.writer = nil
	iw.config = nil
	iw.state = nil
	iw.indented = false
	indentWriterPool.Put(iw)
}

// WriteArrow writes "=> " to the writer.
func WriteArrow(w io.Writer) {
	_, _ = WriteString(w, "=> ")
}

func writeArrowWrappedString(w io.Writer, s string) {
	WriteArrow(w)
	_, _ = WriteString(w, s)
	WriteArrow(w)
}

// WriteNil writes "<nil>" to the writer.
func WriteNil(w io.Writer) {
	_, _ = WriteString(w, "<nil>")
}

func writeTruncated(w io.Writer) {
	_, _ = WriteString(w, "<truncated>")
}

func writeLenCapReflect(w io.Writer, v reflect.Value) {
	writeLenCap(w, v.Len(), v.Cap())
}

func writeLenCap(w io.Writer, ln int, cp int) {
	_, _ = WriteString(w, "(len=")
	_, _ = strconvio.WriteInt(w, int64(ln), 10)
	_, _ = WriteString(w, " cap=")
	_, _ = strconvio.WriteInt(w, int64(cp), 10)
	_, _ = WriteString(w, ")")
}

// WriteString writes a string to the writer.
func WriteString(w io.Writer, s string) (int, error) {
	return writeString(w, s)
}
