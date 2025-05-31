package pretty

import (
	"reflect"

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
	if !vw.match(v) {
		return false
	}
	if !v.CanInterface() {
		return false
	}
	i := v.Interface().(interface{ Text(base int) string }) //nolint:forcetypeassert // Check above.
	text := i.Text(vw.Base)
	write.MustString(st.Writer, text)
	return true
}

func (vw *MathBigIntWriter) match(v reflect.Value) bool {
	kind := v.Kind()
	if kind != reflect.Pointer {
		return false
	}
	elemV := v.Elem()
	elemKind := elemV.Kind()
	if elemKind != reflect.Struct {
		return false
	}
	elemTyp := elemV.Type()
	return elemTyp.PkgPath() == "math/big" && elemTyp.Name() == "Int"
}
