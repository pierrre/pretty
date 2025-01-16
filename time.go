package pretty

import (
	"reflect"
	"time"

	"github.com/pierrre/pretty/internal/write"
)

type TimeValueWriter struct {
	Format string
}

func NewTimeValueWriter() *TimeValueWriter {
	return &TimeValueWriter{
		Format: time.RFC3339Nano,
	}
}

var timeType = reflect.TypeFor[time.Time]()

func (wv *TimeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Type() != timeType {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	tm := v.Interface().(time.Time) //nolint:forcetypeassert // Check above.
	bp := bytesPool.Get()
	defer bytesPool.Put(bp)
	*bp = tm.AppendFormat((*bp)[:0], wv.Format)
	write.Must(st.Writer.Write(*bp))
	return true
}
