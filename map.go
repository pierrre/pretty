package pretty

import (
	"cmp"
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/pierrre/go-libs/syncutil"
)

// MapValueWriter is a [ValueWriter] that handles map values.
//
// It should be created with [NewMapValueWriter].
type MapValueWriter struct {
	ValueWriter
	// ShowLen shows the len.
	// Default: true.
	ShowLen bool
	// ShowAddr shows the address.
	// Default: false.
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
		ShowAddr:           false,
		SortKeys:           false,
		SortKeysCmpDefault: mapSortKeysCmpDefault,
		MaxLen:             0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MapValueWriter) WriteValue(w io.Writer, st State, v reflect.Value) bool {
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
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(w)
	writeString(w, "{")
	if v.Len() > 0 {
		writeString(w, "\n")
		st.IndentLevel++
		if vw.SortKeys {
			vw.writeSorted(w, st, v)
		} else {
			vw.writeUnsorted(w, st, v)
		}
		st.IndentLevel--
		st.writeIndent(w)
	}
	writeString(w, "}")
	return true
}

func (vw *MapValueWriter) writeSorted(w io.Writer, st State, v reflect.Value) {
	keys := vw.getSortedKeys(v)
	for i, key := range keys {
		ok := vw.writeEntry(w, st, key, v.MapIndex(key), i)
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

func (vw *MapValueWriter) writeUnsorted(w io.Writer, st State, v reflect.Value) {
	if v.CanInterface() {
		vw.writeUnsortedExported(w, st, v)
	} else {
		vw.writeUnsortedUnexported(w, st, v)
	}
}

var reflectValuePools = syncutil.MapFor[reflect.Type, *syncutil.PoolFor[*reflect.Value]]{}

func getReflectValuePool(typ reflect.Type) *syncutil.PoolFor[*reflect.Value] {
	pool, ok := reflectValuePools.Load(typ)
	if ok {
		return pool
	}
	pool = &syncutil.PoolFor[*reflect.Value]{
		New: func() *reflect.Value {
			v := reflect.New(typ).Elem()
			return &v
		},
	}
	reflectValuePools.Store(typ, pool)
	return pool
}

func (vw *MapValueWriter) writeUnsortedExported(w io.Writer, st State, v reflect.Value) {
	iter := v.MapRange()
	typ := v.Type()
	keyPool := getReflectValuePool(typ.Key())
	valuePool := getReflectValuePool(typ.Elem())
	keyP := keyPool.Get()
	valueP := valuePool.Get()
	key := *keyP
	value := *valueP
	defer func() {
		key.SetZero()
		value.SetZero()
		keyPool.Put(keyP)
		valuePool.Put(valueP)
	}()
	for i := 0; iter.Next(); i++ {
		key.SetIterKey(iter)
		value.SetIterValue(iter)
		ok := vw.writeEntry(w, st, key, value, i)
		if !ok {
			break
		}
	}
}

func (vw *MapValueWriter) writeUnsortedUnexported(w io.Writer, st State, v reflect.Value) {
	iter := v.MapRange()
	for i := 0; iter.Next(); i++ {
		key := iter.Key()
		value := iter.Value()
		ok := vw.writeEntry(w, st, key, value, i)
		if !ok {
			break
		}
	}
}

func (vw *MapValueWriter) writeEntry(w io.Writer, st State, key reflect.Value, value reflect.Value, i int) bool {
	st.writeIndent(w)
	if vw.MaxLen > 0 && i >= vw.MaxLen {
		writeTruncated(w)
		writeString(w, "\n")
		return false
	}
	mustHandle(vw.ValueWriter(w, st, key))
	writeString(w, ": ")
	mustHandle(vw.ValueWriter(w, st, value))
	writeString(w, ",\n")
	return true
}
