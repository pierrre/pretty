// Package pretty provides utilities to pretty print values.
package pretty

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"
	"unsafe" //nolint:depguard // Required for string to []byte conversion.

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/strconvio"
	"golang.org/x/exp/slices" // TODO use stdlib when Go 1.21 is available.
)

// Write writes the value to the writer with DefaultConfig.
func Write(w io.Writer, vi any) {
	DefaultConfig.Write(w, vi)
}

// String returns the value as a string with DefaultConfig.
func String(vi any) string {
	return DefaultConfig.String(vi)
}

// Formatter returns a fmt.Formatter for the value with DefaultConfig.
func Formatter(vi any) fmt.Formatter {
	return DefaultConfig.Formatter(vi)
}

// DefaultConfig is the default config.
var DefaultConfig = NewConfig()

// Config is a configuration used to pretty print values.
//
// It should be created with NewConfig().
type Config struct {
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
	// ValueWriters is the list of ValueWriter used to write values.
	// Default: reflect.Value, error, []byte, interface{ Bytes() []byte }, fmt.Stringer.
	ValueWriters []ValueWriter
}

// NewConfig creates a new Config initialized with default values.
func NewConfig() *Config {
	return &Config{
		Indent: "\t",
		ValueWriters: []ValueWriter{
			NewReflectValueValueWriter(),
			NewErrorValueWriter(),
			NewBytesValueWriter(0),
			NewByteserValueWriter(0),
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
	c.writeTypeAndValue(w, st, v)
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

// Formatter returns a fmt.Formatter for the value.
func (c *Config) Formatter(v any) fmt.Formatter {
	return &formatter{
		config: c,
		value:  v,
	}
}

var (
	indentCacheLock sync.Mutex
	indentCache     = map[string][]byte{}
)

func (c *Config) writeIndent(w io.Writer, st *State) {
	if st.Indent <= 0 {
		return
	}
	if st.Indent == 1 {
		_, _ = writeString(w, c.Indent)
		return
	}
	indentCacheLock.Lock()
	b := indentCache[c.Indent]
	l := len(c.Indent) * st.Indent
	if len(b) < l {
		b = bytes.Repeat([]byte(c.Indent), st.Indent)
		indentCache[c.Indent] = b
	}
	indentCacheLock.Unlock()
	_, _ = w.Write(b[:l])
}

func (c *Config) checkRecursion(w io.Writer, st *State, v reflect.Value) bool {
	vp := v.Pointer()
	if slices.Contains(st.Pointers, vp) {
		_, _ = writeString(w, "<recursion>")
		return true
	}
	st.Pointers = append(st.Pointers, vp)
	return false
}

func (c *Config) endRecursion(st *State) {
	st.Pointers = st.Pointers[:len(st.Pointers)-1]
}

func (c *Config) writeTypeAndValue(w io.Writer, st *State, v reflect.Value) {
	if !v.IsValid() {
		writeNil(w)
		return
	}
	if v.Kind() == reflect.Interface {
		c.writeTypeAndValue(w, st, v.Elem())
		return
	}
	if c.MaxDepth > 0 && st.Depth >= c.MaxDepth {
		_, _ = writeString(w, "<max depth>")
		return
	}
	st.Depth++
	_, _ = writeString(w, "(")
	c.writeType(w, v.Type())
	_, _ = writeString(w, ") ")
	c.writeValue(w, st, v)
	st.Depth--
}

func (c *Config) writeType(w io.Writer, typ reflect.Type) {
	_, _ = writeString(w, typ.String())
}

func (c *Config) writeValue(w io.Writer, st *State, v reflect.Value) {
	if c.writeValueWithValueWriter(w, st, v) {
		return
	}
	c.writeValueDefault(w, st, v)
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

//nolint:gocyclo // We need to handle all kinds.
func (c *Config) writeValueDefault(w io.Writer, st *State, v reflect.Value) {
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
	_, _ = writeString(w, "0x")
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
	_, _ = writeString(w, "(len=")
	_, _ = strconvio.WriteInt(w, int64(len(s)), 10)
	_, _ = writeString(w, ") ")
	truncated := false
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
		truncated = true
	}
	_, _ = strconvio.WriteQuote(w, s)
	if truncated {
		_, _ = writeString(w, " ")
		writeTruncated(w)
	}
}

func (c *Config) writeChan(w io.Writer, v reflect.Value) {
	if v.IsNil() {
		writeNil(w)
		return
	}
	writeLenCapReflect(w, v)
}

func (c *Config) writeFunc(w io.Writer, v reflect.Value) {
	if v.IsNil() {
		writeNil(w)
		return
	}
	name := runtime.FuncForPC(v.Pointer()).Name()
	_, _ = writeString(w, name)
}

func (c *Config) writePointer(w io.Writer, st *State, v reflect.Value) {
	if c.checkRecursion(w, st, v) {
		return
	}
	writeArrow(w)
	c.writeTypeAndValue(w, st, v.Elem())
	c.endRecursion(st)
}

func (c *Config) writeUnsafePointer(w io.Writer, v reflect.Value) {
	_, _ = writeString(w, "0x")
	_, _ = strconvio.WriteUint(w, uint64(uintptr(v.UnsafePointer())), 16)
}

func (c *Config) writeArray(w io.Writer, st *State, v reflect.Value) {
	l := v.Len()
	truncated := false
	if c.SliceMaxLen > 0 && l > c.SliceMaxLen {
		l = c.SliceMaxLen
		truncated = true
	}
	_, _ = writeString(w, "{\n")
	if v.Len() > 0 {
		st.Indent++
		for i := 0; i < l; i++ {
			c.writeIndent(w, st)
			c.writeTypeAndValue(w, st, v.Index(i))
			_, _ = writeString(w, ",\n")
		}
		if truncated {
			c.writeIndent(w, st)
			writeTruncated(w)
			_, _ = writeString(w, "\n")
		}
		st.Indent--
	}
	c.writeIndent(w, st)
	_, _ = writeString(w, "}")
}

func (c *Config) writeSlice(w io.Writer, st *State, v reflect.Value) {
	if v.IsNil() {
		writeNil(w)
		return
	}
	if c.checkRecursion(w, st, v) {
		return
	}
	writeLenCapReflect(w, v)
	_, _ = writeString(w, " ")
	c.writeArray(w, st, v)
	c.endRecursion(st)
}

func (c *Config) writeMap(w io.Writer, st *State, v reflect.Value) {
	if v.IsNil() {
		writeNil(w)
		return
	}
	if c.checkRecursion(w, st, v) {
		return
	}
	_, _ = writeString(w, "(len=")
	_, _ = strconvio.WriteInt(w, int64(v.Len()), 10)
	_, _ = writeString(w, ") {\n")
	if v.Len() > 0 {
		st.Indent++
		if c.MapSortKeys {
			c.writeMapSorted(w, st, v)
		} else {
			c.writeMapUnsorted(w, st, v)
		}
		st.Indent--
	}
	c.writeIndent(w, st)
	_, _ = writeString(w, "}")
	c.endRecursion(st)
}

func (c *Config) writeMapSorted(w io.Writer, st *State, v reflect.Value) {
	keys := v.MapKeys()
	c.sortValues(v.Type().Key(), keys)
	for i, key := range keys {
		ok := c.writeMapEntry(w, st, key, v.MapIndex(key), i)
		if !ok {
			break
		}
	}
}

func (c *Config) sortValues(typ reflect.Type, vs []reflect.Value) {
	less := c.getLess(typ)
	slices.SortFunc(vs, less)
}

func (c *Config) getLess(typ reflect.Type) func(a, b reflect.Value) bool {
	switch typ.Kind() { //nolint:exhaustive // Optimized for common kinds, the default case is less optimized.
	case reflect.Bool:
		return lessBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lessInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return lessUint
	case reflect.Float32, reflect.Float64:
		return lessFloat
	case reflect.String:
		return lessString
	default:
		return c.lessDefault
	}
}

func lessBool(a, b reflect.Value) bool {
	return !a.Bool() && b.Bool()
}

func lessInt(a, b reflect.Value) bool {
	return a.Int() < b.Int()
}

func lessUint(a, b reflect.Value) bool {
	return a.Uint() < b.Uint()
}

func lessFloat(a, b reflect.Value) bool {
	return a.Float() < b.Float()
}

func lessString(a, b reflect.Value) bool {
	return a.String() < b.String()
}

func (c *Config) lessDefault(a, b reflect.Value) bool {
	return c.string(a) < c.string(b)
}

var typeInterface = reflect.TypeOf((*any)(nil)).Elem()

var reflectValuePool = &sync.Pool{
	New: func() any {
		return reflect.New(typeInterface).Elem()
	},
}

func (c *Config) writeMapUnsorted(w io.Writer, st *State, v reflect.Value) {
	iter := v.MapRange()
	keyItf := reflectValuePool.Get()
	valueItf := reflectValuePool.Get()
	key := keyItf.(reflect.Value)     //nolint:forcetypeassert // The pool only contains *State.
	value := valueItf.(reflect.Value) //nolint:forcetypeassert // The pool only contains *State.
	for i := 0; iter.Next(); i++ {
		key.SetIterKey(iter)
		value.SetIterValue(iter)
		ok := c.writeMapEntry(w, st, key, value, i)
		if !ok {
			break
		}
	}
	key.SetZero()
	value.SetZero()
	reflectValuePool.Put(keyItf)
	reflectValuePool.Put(valueItf)
}

func (c *Config) writeMapEntry(w io.Writer, st *State, key reflect.Value, value reflect.Value, i int) bool {
	c.writeIndent(w, st)
	if c.MapMaxLen > 0 && i >= c.MapMaxLen {
		writeTruncated(w)
		_, _ = writeString(w, "\n")
		return false
	}
	c.writeTypeAndValue(w, st, key)
	_, _ = writeString(w, ": ")
	c.writeTypeAndValue(w, st, value)
	_, _ = writeString(w, ",\n")
	return true
}

func (c *Config) writeStruct(w io.Writer, st *State, v reflect.Value) {
	_, _ = writeString(w, "{\n")
	st.Indent++
	names := getStructFieldNames(v.Type())
	for i, name := range names {
		c.writeIndent(w, st)
		_, _ = w.Write(name)
		_, _ = writeString(w, ": ")
		c.writeTypeAndValue(w, st, v.Field(i))
		_, _ = writeString(w, ",\n")
	}
	st.Indent--
	c.writeIndent(w, st)
	_, _ = writeString(w, "}")
}

var (
	structFieldNamesCacheLock sync.Mutex
	structFieldNamesCache     = map[reflect.Type][][]byte{}
)

func getStructFieldNames(typ reflect.Type) [][]byte {
	structFieldNamesCacheLock.Lock()
	names, ok := structFieldNamesCache[typ]
	if !ok {
		names = make([][]byte, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			names[i] = []byte(typ.Field(i).Name)
		}
		structFieldNamesCache[typ] = names
	}
	structFieldNamesCacheLock.Unlock()
	return names
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
	Depth    int
	Indent   int
	Pointers []uintptr
}

func (st *State) reset() {
	st.Depth = 0
	st.Indent = 0
	st.Pointers = st.Pointers[:0]
}

// ValueWriter is a function that writes a value.
// It can be used to override the default behavior.
//
// It returns true if it handled the value, false otherwise.
//
// Implementations must check v.CanInterface() before using v.Interface().
type ValueWriter func(c *Config, w io.Writer, st *State, v reflect.Value) bool

var typeReflectValue = reflect.TypeOf(reflect.Value{})

// NewReflectValueValueWriter returns a ValueWriter that writes reflect.Value.
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
	writeArrow(w)
	c.writeTypeAndValue(w, st, rv)
	return true
}

var typeError = reflect.TypeOf((*error)(nil)).Elem()

// NewErrorValueWriter returns a ValueWriter that writes error.
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
	err := v.Interface().(error) //nolint:forcetypeassert // Checked above.
	_, _ = writeString(w, "=> .Error() => ")
	_, _ = strconvio.WriteQuote(w, err.Error())
	return true
}

var bytesType = reflect.TypeOf([]byte(nil))

// NewBytesValueWriter returns a ValueWriter that writes []byte with encoding/hex.Dumper.
func NewBytesValueWriter(maxLen int) ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		return writeBytes(c, w, st, v, maxLen)
	}
}

func writeBytes(c *Config, w io.Writer, st *State, v reflect.Value, maxLen int) bool {
	if v.Type() != bytesType {
		return false
	}
	if v.IsNil() {
		writeNil(w)
		return true
	}
	writeLenCapReflect(w, v)
	b := v.Bytes()
	writeBytesCommon(c, w, st, b, maxLen)
	return true
}

type byteser interface {
	Bytes() []byte
}

var byteserType = reflect.TypeOf((*byteser)(nil)).Elem()

// NewByteserValueWriter returns a ValueWriter that writes interface { Bytes() []byte } with encoding/hex.Dumper.
func NewByteserValueWriter(maxLen int) ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		return writeByteser(c, w, st, v, maxLen)
	}
}

func writeByteser(c *Config, w io.Writer, st *State, v reflect.Value, maxLen int) bool {
	if !v.CanInterface() {
		return false
	}
	if !v.Type().Implements(byteserType) {
		return false
	}
	br := v.Interface().(byteser) //nolint:forcetypeassert // Checked above.
	b := br.Bytes()
	_, _ = writeString(w, "=> .Bytes() => ")
	if b == nil {
		writeNil(w)
		return true
	}
	writeLenCap(w, len(b), cap(b))
	writeBytesCommon(c, w, st, b, maxLen)
	return true
}

func writeBytesCommon(c *Config, w io.Writer, st *State, b []byte, maxLen int) {
	truncated := false
	if maxLen > 0 && len(b) > maxLen {
		b = b[:maxLen]
		truncated = true
	}
	_, _ = writeString(w, "\n")
	st.Indent++
	iw := GetIndentWriter(w, c, st, false)
	d := hex.Dumper(iw)
	_, _ = d.Write(b)
	_ = d.Close()
	iw.Release()
	if truncated {
		c.writeIndent(w, st)
		writeTruncated(w)
	}
	st.Indent--
}

var typeStringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

// NewStringerValueWriter returns a ValueWriter that writes fmt.Stringer.
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
	sr := v.Interface().(fmt.Stringer) //nolint:forcetypeassert // Checked above.
	s := sr.String()
	_, _ = writeString(w, "=> .String() => ")
	writeStringValue(w, s, maxLen)
	return true
}

// NewFilterValueWriter returns a ValueWriter that calls the provided ValueWriter if f returns true.
//
// It allows to enable/disable a ValueWriter for specific values/types.
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

// NewDefaultValueWriter returns a ValueWriter that writes the value with the default behavior, bypassing all ValueWriters.
//
// It should be used with NewFilterValueWriter() in order to filter specific types.
func NewDefaultValueWriter() ValueWriter {
	return func(c *Config, w io.Writer, st *State, v reflect.Value) bool {
		c.writeValueDefault(w, st, v)
		return true
	}
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
			iw.config.writeIndent(iw.writer, iw.state)
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

func writeArrow(w io.Writer) {
	_, _ = writeString(w, "=> ")
}

func writeNil(w io.Writer) {
	_, _ = writeString(w, "<nil>")
}

func writeTruncated(w io.Writer) {
	_, _ = writeString(w, "<truncated>")
}

func writeLenCapReflect(w io.Writer, v reflect.Value) {
	writeLenCap(w, v.Len(), v.Cap())
}

func writeLenCap(w io.Writer, ln int, cp int) {
	_, _ = writeString(w, "(len=")
	_, _ = strconvio.WriteInt(w, int64(ln), 10)
	_, _ = writeString(w, " cap=")
	_, _ = strconvio.WriteInt(w, int64(cp), 10)
	_, _ = writeString(w, ")")
}

func writeString(w io.Writer, s string) (int, error) {
	return w.Write(unsafeStringToBytes(s)) //nolint:wrapcheck // The error is not wrapped.
}

func unsafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
