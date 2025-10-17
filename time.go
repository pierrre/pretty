package pretty

import (
	"reflect"
	"time"

	"github.com/pierrre/pretty/internal/write"
)

// TimeWriter is a [ValueWriter] that handles [time.Time], [time.Duration] and [time.Location].
//
// It should be created with [NewTimeWriter].
type TimeWriter struct {
	Time     *TimeTimeWriter
	Duration *TimeDurationWriter
	Location *TimeLocationWriter
}

// NewTimeWriter creates a new [TimeWriter].
func NewTimeWriter() *TimeWriter {
	return &TimeWriter{
		Time:     NewTimeTimeWriter(),
		Duration: NewTimeDurationWriter(),
		Location: NewTimeLocationWriter(),
	}
}

// WriteValue implements [ValueWriter].
func (vw *TimeWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.Time != nil && vw.Time.WriteValue(st, v) {
		return true
	}
	if vw.Duration != nil && vw.Duration.WriteValue(st, v) {
		return true
	}
	if vw.Location != nil && vw.Location.WriteValue(st, v) {
		return true
	}
	return false
}

// Supports implements [SupportChecker].
func (vw *TimeWriter) Supports(typ reflect.Type) ValueWriter {
	if w := callSupportCheckerPointer(vw.Time, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Duration, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Location, typ); w != nil {
		return w
	}
	return nil
}

// TimeTimeWriter is a [ValueWriter] that handles [time.Time] values.
//
// It should be created with [NewTimeTimeWriter].
type TimeTimeWriter struct {
	// Format is the format of the time.
	// Default: [time.RFC3339Nano].
	Format string

	// Location to convert the time before formatting.
	// Default: nil (no conversion).
	Location *time.Location
}

// NewTimeTimeWriter creates a new [TimeTimeWriter] with default values.
func NewTimeTimeWriter() *TimeTimeWriter {
	return &TimeTimeWriter{
		Format: time.RFC3339Nano,
	}
}

var timeTimeType = reflect.TypeFor[time.Time]()

// WriteValue implements [ValueWriter].
func (vw *TimeTimeWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Struct || v.Type() != timeTimeType {
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
func (vw *TimeTimeWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == timeTimeType {
		res = vw
	}
	return res
}

// TimeDurationWriter is a [ValueWriter] that handles [time.Duration] values.
//
// It should be created with [NewTimeDurationWriter].
type TimeDurationWriter struct{}

// NewTimeDurationWriter creates a new [TimeDurationWriter].
func NewTimeDurationWriter() *TimeDurationWriter {
	return &TimeDurationWriter{}
}

var timeDurationType = reflect.TypeFor[time.Duration]()

// WriteValue implements [ValueWriter].
func (vw *TimeDurationWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Int64 || v.Type() != timeDurationType {
		return false
	}
	d := time.Duration(v.Int())
	write.MustString(st.Writer, d.String())
	return true
}

// Supports implements [SupportChecker].
func (vw *TimeDurationWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == timeDurationType {
		res = vw
	}
	return res
}

// TimeLocationWriter is a [ValueWriter] that handles [time.Location] values.
//
// It should be created with [NewTimeLocationWriter].
type TimeLocationWriter struct{}

// NewTimeLocationWriter creates a new [TimeLocationWriter].
func NewTimeLocationWriter() *TimeLocationWriter {
	return &TimeLocationWriter{}
}

var timeLocationType = reflect.TypeFor[*time.Location]()

// WriteValue implements [ValueWriter].
func (vw *TimeLocationWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer || v.Type() != timeLocationType {
		return false
	}
	if v.IsNil() {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	loc, _ := reflect.TypeAssert[*time.Location](v)
	write.MustString(st.Writer, loc.String())
	return true
}

// Supports implements [SupportChecker].
func (vw *TimeLocationWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ == timeLocationType {
		res = vw
	}
	return res
}
