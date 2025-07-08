package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/must"
	"github.com/pierrre/pretty/internal/write"
)

// MathBigWriter is a [ValueWriter] thats handles values from the [math/big] package.
//
// It should be created with [NewMathBigWriter].
type MathBigWriter struct {
	Int   *MathBigIntWriter
	Float *MathBigFloatWriter
	Rat   *MathBigRatWriter
}

// WriteValue implements [ValueWriter].
func (vw *MathBigWriter) WriteValue(st *State, v reflect.Value) bool {
	if vw.Int != nil && vw.Int.WriteValue(st, v) {
		return true
	}
	if vw.Float != nil && vw.Float.WriteValue(st, v) {
		return true
	}
	if vw.Rat != nil && vw.Rat.WriteValue(st, v) {
		return true
	}
	return false
}

// Supports implements [SupportChecker].
func (vw *MathBigWriter) Supports(typ reflect.Type) ValueWriter {
	if w := callSupportCheckerPointer(vw.Int, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Float, typ); w != nil {
		return w
	}
	if w := callSupportCheckerPointer(vw.Rat, typ); w != nil {
		return w
	}
	return nil
}

// MathBigIntWriter is a [ValueWriter] that handles [math/big.Int] values.
//
// It should be created with [NewMathBigIntWriter].
type MathBigIntWriter struct {
	// Base is the base used to format the integer.
	// Default: 10.
	Base int
}

// NewMathBigIntWriter creates a new [MathBigIntWriter] with default values.
func NewMathBigIntWriter() *MathBigIntWriter {
	return &MathBigIntWriter{
		Base: 10,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MathBigIntWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	elemTyp := v.Type().Elem()
	if elemTyp.Kind() != reflect.Struct || elemTyp.PkgPath() != "math/big" || elemTyp.Name() != "Int" {
		return false
	}
	i, ok := itfassert.Assert[interface {
		Append(buf []byte, base int) []byte
	}](v)
	if !ok {
		return false
	}
	bp := bytesPool.Get()
	*bp = i.Append((*bp)[:0], vw.Base)
	write.Must(st.Writer.Write(*bp))
	bytesPool.Put(bp)
	return true
}

// Supports implements [SupportChecker].
func (vw *MathBigIntWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Pointer {
		elemTyp := typ.Elem()
		if elemTyp.Kind() == reflect.Struct && elemTyp.PkgPath() == "math/big" && elemTyp.Name() == "Int" {
			res = vw
		}
	}
	return res
}

// NewMathBigWriter creates a new [MathBigWriter] with default values.
func NewMathBigWriter() *MathBigWriter {
	return &MathBigWriter{
		Float: NewMathBigFloatWriter(),
		Int:   NewMathBigIntWriter(),
		Rat:   NewMathBigRatWriter(),
	}
}

// MathBigFloatWriter is a [ValueWriter] that handles [math/big.Float] values.
//
// It should be created with [NewMathBigFloatWriter].
type MathBigFloatWriter struct {
	// Format is the format used to format the float.
	// Default: 'g'.
	Format byte
	// Precision is the precision used to format the float.
	// Default: -1.
	Precision int
}

// NewMathBigFloatWriter creates a new [MathBigFloatWriter] with default values.
func NewMathBigFloatWriter() *MathBigFloatWriter {
	return &MathBigFloatWriter{
		Format:    'g',
		Precision: -1,
	}
}

// WriteValue implements [ValueWriter].
func (vw *MathBigFloatWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	elemTyp := v.Type().Elem()
	if elemTyp.Kind() != reflect.Struct || elemTyp.PkgPath() != "math/big" || elemTyp.Name() != "Float" {
		return false
	}
	i, ok := itfassert.Assert[interface {
		Append(buf []byte, fmt byte, prec int) []byte
	}](v)
	if !ok {
		return false
	}
	bp := bytesPool.Get()
	*bp = i.Append((*bp)[:0], vw.Format, vw.Precision)
	write.Must(st.Writer.Write(*bp))
	bytesPool.Put(bp)
	return true
}

// Supports implements [SupportChecker].
func (vw *MathBigFloatWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Pointer {
		elemTyp := typ.Elem()
		if elemTyp.Kind() == reflect.Struct && elemTyp.PkgPath() == "math/big" && elemTyp.Name() == "Float" {
			res = vw
		}
	}
	return res
}

// MathBigRatWriter is a [ValueWriter] that handles [math/big.Rat] values.
//
// It should be created with [NewMathBigRatWriter].
type MathBigRatWriter struct{}

// NewMathBigRatWriter creates a new [MathBigRatWriter] with default values.
func NewMathBigRatWriter() *MathBigRatWriter {
	return &MathBigRatWriter{}
}

// WriteValue implements [ValueWriter].
func (vw *MathBigRatWriter) WriteValue(st *State, v reflect.Value) bool {
	if v.Kind() != reflect.Pointer {
		return false
	}
	elemTyp := v.Type().Elem()
	if elemTyp.Kind() != reflect.Struct || elemTyp.PkgPath() != "math/big" || elemTyp.Name() != "Rat" {
		return false
	}
	i, ok := itfassert.Assert[interface {
		AppendText(b []byte) ([]byte, error)
	}](v)
	if !ok {
		return false
	}
	bp := bytesPool.Get()
	var err error
	*bp, err = i.AppendText((*bp)[:0])
	must.NoError(err)
	write.Must(st.Writer.Write(*bp))
	bytesPool.Put(bp)
	return true
}

// Supports implements [SupportChecker].
func (vw *MathBigRatWriter) Supports(typ reflect.Type) ValueWriter {
	var res ValueWriter
	if typ.Kind() == reflect.Pointer {
		elemTyp := typ.Elem()
		if elemTyp.Kind() == reflect.Struct && elemTyp.PkgPath() == "math/big" && elemTyp.Name() == "Rat" {
			res = vw
		}
	}
	return res
}
