package pretty

import (
	"reflect"

	"github.com/pierrre/pretty/internal/itfassert"
	"github.com/pierrre/pretty/internal/write"
)

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
	i, ok := itfassert.Assert[interface{ Text(base int) string }](v)
	if !ok {
		return false
	}
	text := i.Text(vw.Base)
	write.MustString(st.Writer, text)
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
