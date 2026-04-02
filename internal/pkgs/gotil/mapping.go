package gotil

import (
	"context"
	"errors"
	"reflect"

	"github.com/spf13/cast"
)

// MappingFunc 映射函数
type MappingFunc[K comparable] func(i int) (key K, ptr any)

type DBLoader[K comparable] interface {
	LoadMany(ctx context.Context, keys []K) (map[K]any, error)
}

// keyPtr key 与 ptr
type keyPtr[K comparable] struct {
	key K
	ptr any
}

// FillToSlice 根据 MappingFunc[K] 收集 key 与 ptr，批量加载后通过 castPtr 写回切片
// keyValid 过滤无效 key，如 int 用 func(k int) bool { return k > 0 }，string 用 func(k string) bool { return len(k) > 0 }
func FillToSlice[K comparable](ctx context.Context, loader DBLoader[K], slice any, f MappingFunc[K], keyValid func(K) bool) error {
	keyPtrs, keys, err := collectKeysAndPtrs(slice, f, keyValid)
	if err != nil || len(keys) == 0 {
		return err
	}

	valueMap, err := loader.LoadMany(ctx, keys)
	if err != nil {
		return err
	}

	return fillValues(keyPtrs, valueMap)
}

// collectKeysAndPtrs 收集有效的 key 和 ptr
func collectKeysAndPtrs[K comparable](slice any, f MappingFunc[K], keyValid func(K) bool) ([]keyPtr[K], []K, error) {
	var keyPtrs []keyPtr[K]
	keySet := make(map[K]bool)
	var keys []K

	// 遍历 slice,获取长度
	sliceLen, err := getSliceLen(slice)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < sliceLen; i++ {
		key, ptr := f(i)
		if !keyValid(key) || ptr == nil {
			continue
		}

		keyPtrs = append(keyPtrs, keyPtr[K]{key: key, ptr: ptr})
		if !keySet[key] {
			keys = append(keys, key)
			keySet[key] = true
		}
	}
	return keyPtrs, keys, nil
}

// fillValues 将加载的数据写回指针
func fillValues[K comparable](keyPtrs []keyPtr[K], valueMap map[K]any) error {
	for _, kp := range keyPtrs {
		if v, ok := valueMap[kp.key]; !ok {
			continue
		} else {
			castPtr(kp.ptr, v)
		}
	}
	return nil
}

// getSliceLen 获取 slice 长度
func getSliceLen(slice any) (l int, err error) {
	v := reflect.ValueOf(slice)
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Slice {
		err = errors.New("invalid mapping slice")
		return
	}
	l = v.Len()
	return
}

// castPtr 把 v 写回 ptr
func castPtr(ptr, v any) {
	switch t := ptr.(type) {
	case *int:
		*t = cast.ToInt(v)
	case *string:
		*t = cast.ToString(v)
	}
}
