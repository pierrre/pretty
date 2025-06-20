package pretty

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// MapWriter is a [ValueWriter] that handles map values.
//
// It should be created with [NewMapWriter].
type MapWriter struct {
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
	// ShowKeysInfos shows keys infos.
	// Default: false.
	ShowKeysInfos bool
	// MaxLen is the maximum length of the map.
	// Default: 0 (no limit).
	MaxLen int
}

// NewMapWriter creates a new [MapWriter] with default values.
func NewMapWriter(vw ValueWriter) *MapWriter {
	return &MapWriter{
		ValueWriter:   vw,
		ShowLen:       true,
		ShowAddr:      false,
		SortKeys:      false,
		ShowKeysInfos: false,
		MaxLen:        0,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MapWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Map {
		return false
	}
	if checkNil(st.Writer, v) {
		return true
	}
	l := v.Len()
	infos{
		showLen:  vw.ShowLen,
		len:      l,
		showAddr: vw.ShowAddr,
		addr:     uintptr(v.UnsafePointer()),
	}.writeWithTrailingSpace(st)
	write.MustString(st.Writer, "{")
	if l > 0 {
		write.MustString(st.Writer, "\n")
		st.IndentLevel++
		if vw.SortKeys {
			vw.writeSorted(st, v)
		} else {
			vw.writeUnsorted(st, v)
		}
		st.IndentLevel--
		st.WriteIndent()
	}
	write.MustString(st.Writer, "}")
	return true
}

func (vw *MapWriter) writeSorted(st *State, v reflect.Value) {
	es := reflectutil.GetSortedMap(v)
	for i, e := range es {
		ok := vw.writeEntry(st, e.Key, e.Value, i)
		if !ok {
			break
		}
	}
	es.Release()
}

func (vw *MapWriter) writeUnsorted(st *State, v reflect.Value) {
	if v.CanInterface() {
		vw.writeUnsortedExported(st, v)
	} else {
		vw.writeUnsortedUnexported(st, v)
	}
}

var reflectValuePools syncutil.Map[reflect.Type, *syncutil.Pool[*reflect.Value]]

func getReflectValuePool(typ reflect.Type) *syncutil.Pool[*reflect.Value] {
	pool, ok := reflectValuePools.Load(typ)
	if ok {
		return pool
	}
	pool = &syncutil.Pool[*reflect.Value]{
		New: func() *reflect.Value {
			v := reflect.New(typ).Elem()
			return &v
		},
	}
	pool, _ = reflectValuePools.LoadOrStore(typ, pool)
	return pool
}

func (vw *MapWriter) writeUnsortedExported(st *State, v reflect.Value) {
	iter := v.MapRange()
	typ := v.Type()
	keyPool := getReflectValuePool(typ.Key())
	valuePool := getReflectValuePool(typ.Elem())
	keyP := keyPool.Get()
	valueP := valuePool.Get()
	key := *keyP
	value := *valueP
	for i := 0; iter.Next(); i++ {
		key.SetIterKey(iter)
		value.SetIterValue(iter)
		ok := vw.writeEntry(st, key, value, i)
		if !ok {
			break
		}
	}
	key.SetZero()
	value.SetZero()
	keyPool.Put(keyP)
	valuePool.Put(valueP)
}

func (vw *MapWriter) writeUnsortedUnexported(st *State, v reflect.Value) {
	iter := v.MapRange()
	for i := 0; iter.Next(); i++ {
		key := iter.Key()
		value := iter.Value()
		ok := vw.writeEntry(st, key, value, i)
		if !ok {
			break
		}
	}
}

func (vw *MapWriter) writeEntry(st *State, key reflect.Value, value reflect.Value, i int) bool {
	st.WriteIndent()
	if vw.MaxLen > 0 && i >= vw.MaxLen {
		writeTruncated(st.Writer)
		write.MustString(st.Writer, "\n")
		return false
	}
	showInfos := st.ShowInfos
	st.ShowInfos = vw.ShowKeysInfos
	must.Handle(vw.ValueWriter.WriteValue(st, key))
	st.ShowInfos = showInfos
	write.MustString(st.Writer, ": ")
	must.Handle(vw.ValueWriter.WriteValue(st, value))
	write.MustString(st.Writer, ",\n")
	return true
}

// Supports implements [SupportChecker].
func (vw *MapWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Map {
		res = vw
	}
	return res
}
