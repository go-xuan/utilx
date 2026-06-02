package slicex

// Contains 是否包含
func Contains[T comparable](slice []T, v T) bool {
	for _, item := range slice {
		if item == v {
			return true
		}
	}
	return false
}

// ContainsAny 是否包含任意元素
func ContainsAny[T comparable](slice []T, args ...T) bool {
	if len(args) == 0 {
		return false
	}
	m := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		m[item] = struct{}{}
	}
	for _, arg := range args {
		if _, ok := m[arg]; ok {
			return true
		}
	}
	return false
}

// ContainsAll 是否包含全部元素
func ContainsAll[T comparable](slice []T, args ...T) bool {
	if len(args) == 0 {
		return true
	}
	m := make(map[T]struct{}, len(slice))
	for _, k := range slice {
		m[k] = struct{}{}
	}
	for _, arg := range args {
		if _, ok := m[arg]; !ok {
			return false
		}
	}
	return true
}

// Distinct 去重多个切片（保持首次出现顺序）
func Distinct[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}
	seen := make(map[T]struct{})
	var result []T
	for _, slice := range slices {
		for _, k := range slice {
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				result = append(result, k)
			}
		}
	}
	return result
}

// Intersect 取交集
func Intersect[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}
	m := make(map[T]int)
	for _, slice := range slices {
		seen := make(map[T]struct{})
		for _, k := range slice {
			if _, ok := seen[k]; !ok {
				m[k]++
				seen[k] = struct{}{}
			}
		}
	}
	var result []T
	for k, count := range m {
		if count == len(slices) {
			result = append(result, k)
		}
	}
	return result
}

// Exclude 移除指定元素
func Exclude[T comparable](target []T, exclude []T) []T {
	if len(target) == 0 || len(exclude) == 0 {
		return target
	}
	m := make(map[T]struct{}, len(exclude))
	for _, k := range exclude {
		m[k] = struct{}{}
	}
	result := make([]T, 0, len(target))
	for _, k := range target {
		if _, ok := m[k]; !ok {
			result = append(result, k)
		}
	}
	return result
}

// Difference 取差集（在 target 中但不在 exclude 中）
func Difference[T comparable](target []T, exclude []T) []T {
	return Exclude(target, exclude)
}

// Union 取并集
func Union[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}
	m := make(map[T]struct{})
	for _, slice := range slices {
		for _, k := range slice {
			m[k] = struct{}{}
		}
	}
	result := make([]T, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// IndexOf 查找元素下标
func IndexOf[T comparable](slice []T, v T) int {
	for i, item := range slice {
		if item == v {
			return i
		}
	}
	return -1
}

// LastIndexOf 查找元素最后出现的下标
func LastIndexOf[T comparable](slice []T, v T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == v {
			return i
		}
	}
	return -1
}

// Keys 获取 map 的所有键
func MapKeys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values 获取 map 的所有值
func MapValues[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}
