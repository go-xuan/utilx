// Package slicex 提供切片操作工具
package slicex

// Value 定义具有主键的接口
type Value interface {
	GetPrimaryKey() string // 获取唯一主键
}

// Conv2Map 数组转 map
func Conv2Map[V Value](slice []V) map[string]V {
	if len(slice) == 0 {
		return nil
	}
	m := make(map[string]V, len(slice))
	for _, item := range slice {
		key := item.GetPrimaryKey()
		m[key] = item
	}
	return m
}

// Conv2Groups 数组分组
func Conv2Groups[V Value](slice []V) map[string][]V {
	if len(slice) == 0 {
		return nil
	}
	groups := make(map[string][]V)
	for _, item := range slice {
		key := item.GetPrimaryKey()
		groups[key] = append(groups[key], item)
	}
	return groups
}

// Conv2MapBy 根据键函数将切片转换为 map
func Conv2MapBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	if len(slice) == 0 {
		return nil
	}
	m := make(map[K]T, len(slice))
	for _, item := range slice {
		m[keyFunc(item)] = item
	}
	return m
}

// Conv2GroupsBy 根据键函数将切片分组
func Conv2GroupsBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	if len(slice) == 0 {
		return nil
	}
	groups := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// Map 对切片中每个元素应用转换函数
func Map[T, U any](slice []T, fn func(T) U) []U {
	if len(slice) == 0 {
		return nil
	}
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// Filter 过滤切片中满足条件的元素
func Filter[T any](slice []T, fn func(T) bool) []T {
	if len(slice) == 0 {
		return nil
	}
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce 归并切片元素为单个值
func Reduce[T any, R any](slice []T, fn func(R, T) R, init R) R {
	result := init
	for _, item := range slice {
		result = fn(result, item)
	}
	return result
}

// FlatMap 对切片中每个元素应用转换函数并展平结果
func FlatMap[T, U any](slice []T, fn func(T) []U) []U {
	if len(slice) == 0 {
		return nil
	}
	// 两次遍历：先算总长度，再一次性分配，避免多次扩容
	subs := make([][]U, len(slice))
	total := 0
	for i, item := range slice {
		subs[i] = fn(item)
		total += len(subs[i])
	}
	result := make([]U, 0, total)
	for _, sub := range subs {
		result = append(result, sub...)
	}
	return result
}

// Chunk 将切片分割为固定大小的块
func Chunk[T any](slice []T, size int) [][]T {
	if len(slice) == 0 || size <= 0 {
		return nil
	}
	result := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result
}

// Reverse 反转切片
func Reverse[T any](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	result := make([]T, len(slice))
	for i, j := len(slice)-1, 0; i >= 0; i, j = i-1, j+1 {
		result[j] = slice[i]
	}
	return result
}

// Unique 去重（保持原有顺序）
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
