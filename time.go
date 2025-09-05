package pretty

import (
	"reflect"
	"time"

	"github.com/pierrre/pretty/internal/write"
)

// TimeWriter is a [ValueWriter] that handles [time.Time] values.
//
// It should be created with [NewTimeWriter].
type TimeWriter struct {
	// Format is the format of the time.
	// Default: [time.RFC3339Nano].
	Format string

	// Location to convert the time before formatting.
	// Default: nil (no conversion).
	Location *time.Location
}

// NewTimeWriter creates a new [TimeWriter] with default values.
func NewTimeWriter() *TimeWriter {
	return &TimeWriter{
		Format: time.RFC3339Nano,
	}
}

var timeType = reflect.TypeFor[time.Time]()

// WriteValue implements [ValueWriter].
func (vw *TimeWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct || v.Type() != timeType {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	tm, _ := reflect.TypeAssert[time.Time](v)
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
func (vw *TimeWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == timeType {
		res = vw
	}
	return res
}
