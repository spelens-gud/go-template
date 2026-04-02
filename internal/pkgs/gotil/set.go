package gotil

// Set 基于 map[T]struct{} 的泛型集合，用于去重与成员判断
type Set[T comparable] map[T]struct{}

// NewSet 创建容量为 capSize 的空集合
func NewSet[T comparable](capSize int) Set[T] {
	return make(Set[T], capSize)
}

// NewSetWithSlice 从切片构建集合,带去重功
func NewSetWithSlice[T comparable](s []T) Set[T] {
	return NewSet[T](len(s)).AddMany(s)
}

// Add 添加元素，返回是否为新加入
func (set Set[T]) Add(v T) (success bool) {
	_, exist := set[v]
	if !exist {
		success = true
		set[v] = struct{}{}
	}
	return
}

// AddMany 批量添加，支持链式调用
func (set Set[T]) AddMany(multi []T) Set[T] {
	for _, v := range multi {
		set.Add(v)
	}
	return set
}

// Slice 转为切片, 遍历顺序不保证
func (set Set[T]) Slice() []T {
	if len(set) == 0 {
		return nil
	}
	ret := make([]T, 0, len(set))
	for v := range set {
		ret = append(ret, v)
	}
	return ret
}

// CheckSliceIn 判断 small 中所有元素是否都在 big 中,按集合包含
func CheckSliceIn[T comparable](small, big []T) (contain bool) {
	if len(small) == 0 {
		return true
	}
	bm := make(map[T]bool, len(big))
	for _, v := range big {
		bm[v] = true
	}
	for _, v := range small {
		if !bm[v] {
			return false
		}
	}
	return true
}

// CheckSetEqual 判断两个切片作为集合是否相等,忽略顺序与重复
func CheckSetEqual[T comparable](a, b []T) (equal bool) {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	am := make(map[T]struct{}, len(a))
	bm := make(map[T]struct{}, len(b))
	for _, v := range a {
		am[v] = struct{}{}
	}
	for _, v := range b {
		if _, has := am[v]; !has {
			return
		}
		bm[v] = struct{}{}
	}
	for _, v := range a {
		if _, has := bm[v]; !has {
			return
		}
	}
	return true
}
