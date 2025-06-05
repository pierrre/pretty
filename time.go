package pretty

import (
	"reflect"
	"time"

	"github.com/pierrre/pretty/internal/write"
)

// TimeValueWriter is a [ValueWriter] that handles [time.Time] values.
//
// It should be created with [NewTimeValueWriter].
type TimeValueWriter struct {
	// Format is the format of the time.
	// Default: [time.RFC3339Nano].
	Format string

	// Location to convert the time before formatting.
	// Default: nil (no conversion).
	Location *time.Location
}

// NewTimeValueWriter creates a new [TimeValueWriter] with default values.
func NewTimeValueWriter() *TimeValueWriter {
	return &TimeValueWriter{
		Format: time.RFC3339Nano,
	}
}

var timeType = reflect.TypeFor[time.Time]()

// WriteValue implements [ValueWriter].
func (vw *TimeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct || v.Type() != timeType {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	tm := v.Interface().(time.Time) //nolint:forcetypeassert // Check above.
	if vw.Location != nil {
		tm = tm.In(vw.Location)
	}
	bp := bytesPool.Get()
	*bp = tm.AppendFormat((*bp)[:0], vw.Format)
	write.Must(st.Writer.Write(*bp))
	bytesPool.Put(bp)
	return true
}

// Supports implements [SupportChecker].
func (vw *TimeValueWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == timeType {
		res = vw
	}
	return res
}
