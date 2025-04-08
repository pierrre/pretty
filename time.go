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
func (wv *TimeValueWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Type() != timeType {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	tm := v.Interface().(time.Time) //nolint:forcetypeassert // Check above.
	if wv.Location != nil {
		tm = tm.In(wv.Location)
	}
	bp := bytesPool.Get()
	defer bytesPool.Put(bp)
	*bp = tm.AppendFormat((*bp)[:0], wv.Format)
	write.Must(st.Writer.Write(*bp))
	return true
}
