package pretty

import (
	"reflect"
	"testing"

	"github.com/pierrre/go-libs/syncutil"
)

func TestGetReflectValuePoolEvictsAtMaxCount(t *testing.T) {
	reflectValuePools.Clear()
	reflectValuePoolsCount.Store(0)
	intTyp := reflect.TypeFor[int]()
	typs := make([]reflect.Type, reflectValuePoolsMaxCount+1)
	for i := range typs {
		typs[i] = reflect.ArrayOf(i+1, intTyp)
	}
	for _, typ := range typs {
		_ = getReflectValuePool(typ)
	}
	count := 0
	reflectValuePools.Range(func(_ reflect.Type, _ *syncutil.Pool[*reflect.Value]) bool {
		count++
		return true
	})
	if count >= len(typs) {
		t.Fatalf("expected pools to be evicted, got count = %d (inserted %d)", count, len(typs))
	}
}
